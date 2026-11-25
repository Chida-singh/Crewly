import type {
  ClientInput,
  ClientProfile,
  Photo,
  Profile,
  Role,
  TalentInput,
  TalentProfile,
  User,
} from "./api";

// Local development switch. `frontend/.env.development` turns it on for `bun run dev`;
// production builds leave it off, so the real OTP flow is the only way in.
export const AUTH_BYPASS = import.meta.env.VITE_AUTH_BYPASS === "true";

// Change to "client" to preview the hiring side of the sheet.
const DEV_ROLE: Role = "talent";

export const DEV_USER: User = {
  id: "dev-user",
  phone: "+919000000000",
  role: DEV_ROLE,
  is_verified: true,
};

// Placeholder sheet so the demo has something to show. Frames live in public/demo.
const DEMO_TALENT: TalentProfile = {
  role: "talent",
  talent_type: "model",
  display_name: "Aditi Rao",
  city: "Bengaluru",
  age: 24,
  height_cm: 173,
  bio: "Editorial and campaign model. Comfortable on location, quick with direction.",
  rate_min: 18000,
  rate_max: 32000,
  photos: [
    { id: "demo-01", url: "/demo/frame-01.jpg", position: 0 },
    { id: "demo-02", url: "/demo/frame-02.jpg", position: 1 },
  ],
};

const DEMO_CLIENT: ClientProfile = {
  role: "client",
  display_name: "Meera Iyer",
  org_name: "Studio Nine, Bengaluru",
};

// Stand-in store so the UI works without the API. Memory only: a refresh resets to the demo.
let profile: Profile | null = { talent: DEMO_TALENT, client: DEMO_CLIENT }[DEV_ROLE];

export function devGetProfile(): Profile | null {
  return profile;
}

export function devSaveTalentProfile(input: TalentInput): TalentProfile {
  const photos = profile?.role === "talent" ? profile.photos : [];
  profile = { role: "talent", ...input, photos };
  return profile;
}

export function devSaveClientProfile(input: ClientInput): ClientProfile {
  profile = { role: "client", ...input };
  return profile;
}

export function devUploadPhoto(file: File): Photo {
  const talent = profile?.role === "talent" ? profile : null;
  if (!talent) throw new Error("Save your sheet before adding frames.");

  const photo: Photo = {
    id: `dev-${Date.now()}`,
    url: URL.createObjectURL(file),
    position: talent.photos.length,
  };
  profile = { ...talent, photos: [...talent.photos, photo] };
  return photo;
}
