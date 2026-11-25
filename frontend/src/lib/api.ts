export const API_BASE = import.meta.env.VITE_API_BASE_URL ?? "";

export type Role = "talent" | "client";

export type User = {
  id: string;
  phone: string;
  role: Role;
  is_verified: boolean;
};

export type Photo = {
  id: string;
  url: string;
  position: number;
};

export type TalentProfile = {
  role: "talent";
  talent_type: string;
  display_name: string;
  city: string;
  age: number;
  height_cm: number | null;
  bio: string;
  rate_min: number;
  rate_max: number;
  photos: Photo[];
};

export type ClientProfile = {
  role: "client";
  display_name: string;
  org_name: string;
};

export type Profile = TalentProfile | ClientProfile;

export type TalentInput = Omit<TalentProfile, "role" | "photos">;

export type ClientInput = Omit<ClientProfile, "role">;

export class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

type CallOptions = {
  method?: string;
  body?: unknown;
  token?: string;
};

async function call<T>(path: string, { method = "GET", body, token }: CallOptions = {}): Promise<T> {
  const headers: Record<string, string> = {};
  if (body !== undefined) headers["Content-Type"] = "application/json";
  if (token) headers.Authorization = `Bearer ${token}`;

  let response: Response;
  try {
    response = await fetch(API_BASE + path, {
      method,
      headers,
      body: body === undefined ? undefined : JSON.stringify(body),
    });
  } catch {
    throw new ApiError(0, "Can't reach the API. Start the backend on port 8080 and try again.");
  }

  if (!response.ok) {
    let message = `Request failed (${response.status}).`;
    try {
      const data = (await response.json()) as { error?: string };
      if (data?.error) message = data.error;
    } catch {
      // keep the fallback message
    }
    throw new ApiError(response.status, message);
  }

  if (response.status === 204) return undefined as T;
  return (await response.json()) as T;
}

export function requestOtp(phone: string): Promise<{ message: string }> {
  return call("/auth/otp", { method: "POST", body: { phone } });
}

export function verifyOtp(
  phone: string,
  code: string,
  role: Role,
): Promise<{ token: string; user: User }> {
  return call("/auth/verify", { method: "POST", body: { phone, code, role } });
}

export async function getProfile(token: string): Promise<Profile | null> {
  try {
    return await call<Profile>("/me/profile", { token });
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) return null;
    throw error;
  }
}

export function saveTalentProfile(token: string, input: TalentInput): Promise<TalentProfile> {
  return call("/me/profile", { method: "PUT", token, body: input });
}

export function saveClientProfile(token: string, input: ClientInput): Promise<ClientProfile> {
  return call("/me/profile", { method: "PUT", token, body: input });
}

export async function uploadPhoto(token: string, file: File): Promise<Photo> {
  const contentType = file.type || "image/jpeg";
  const reserved = await call<{ photo: Photo; upload_url: string }>("/me/photos", {
    method: "POST",
    token,
    body: { content_type: contentType },
  });

  const response = await fetch(assetUrl(reserved.upload_url), {
    method: "PUT",
    headers: { "Content-Type": contentType },
    body: file,
  }).catch(() => {
    throw new ApiError(0, "Upload failed. Check that the storage endpoint is reachable.");
  });

  if (!response.ok) {
    throw new ApiError(response.status, "Upload failed. The file was not stored.");
  }
  return reserved.photo;
}

// Storage URLs may point at the API origin. Route them through the dev proxy instead.
export function assetUrl(url: string): string {
  try {
    const parsed = new URL(url, window.location.origin);
    if (parsed.pathname.startsWith("/uploads/")) return parsed.pathname + parsed.search;
  } catch {
    // not a URL, keep as-is
  }
  return url;
}

export type Health = { status: string; database: string };

export async function getHealth(): Promise<Health> {
  try {
    const response = await fetch(API_BASE + "/health");
    const data = (await response.json()) as Partial<Health>;
    return { status: data.status ?? "unknown", database: data.database ?? "unknown" };
  } catch {
    return { status: "offline", database: "unknown" };
  }
}

export function errorMessage(error: unknown): string {
  if (error instanceof Error && error.message) return error.message;
  return "Something went wrong. Try again.";
}
