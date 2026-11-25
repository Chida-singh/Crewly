import { useCallback, useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";

import { ApiStatus } from "../components/ApiStatus";
import { Field } from "../components/Field";
import {
  assetUrl,
  errorMessage,
  getProfile,
  saveClientProfile,
  saveTalentProfile,
  uploadPhoto,
  type ClientProfile,
  type Profile,
  type TalentProfile,
} from "../lib/api";
import { DEV_USER } from "../lib/dev";
import { useSession } from "../lib/session";

const MAX_PHOTOS = 12;

export function SheetPage() {
  const { session, logout } = useSession();
  const navigate = useNavigate();
  const token = session?.token ?? "dev";
  const user = session?.user ?? DEV_USER;

  const [profile, setProfile] = useState<Profile | null>(null);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoadError(null);
    try {
      setProfile(await getProfile(token));
    } catch (err) {
      setLoadError(errorMessage(err));
    } finally {
      setLoading(false);
    }
  }, [token]);

  useEffect(() => {
    void load();
  }, [load]);

  function signOut() {
    logout();
    navigate("/", { replace: true });
  }

  const displayName = profile?.display_name || "Call sheet";

  return (
    <div className="page">
      <header className="topbar">
        <Link to="/sheet" className="wordmark" aria-label="Crewly sheet">
          CREWLY<span aria-hidden="true">.</span>
        </Link>
        <div className="topbar-right">
          <ApiStatus />
          {session ? null : <span className="dev-badge">Auth off</span>}
          <span className="role-chip">{user.role === "talent" ? "Talent" : "Client"}</span>
          <button type="button" className="link" onClick={signOut}>
            {session ? "Log out" : "Exit demo"}
          </button>
        </div>
      </header>

      <main className="sheet">
        <div className="sheet-head">
          <h1 className="sheet-title">{displayName}</h1>
          <p className="sheet-sub">
            +91 {user.phone.replace(/^\+91/, "")} · {user.role === "talent" ? "Model" : "Client"} ·
            Phone verified
          </p>
        </div>

        {loading ? (
          <div className="sheet-grid" aria-busy="true">
            <div className="sheet-form">
              <div className="skeleton" />
              <div className="skeleton short" />
              <div className="skeleton" />
            </div>
            <div className="sheet-side">
              <div className="skeleton tall" />
            </div>
          </div>
        ) : loadError ? (
          <div className="banner" role="alert">
            <p>{loadError}</p>
            <button type="button" className="btn small" onClick={() => void load()}>
              Try again
            </button>
          </div>
        ) : (
          <div className={`sheet-grid ${user.role === "client" ? "single" : ""}`}>
            <section className="sheet-form" aria-label="Sheet details">
              {user.role === "talent" ? (
                <TalentEditor
                  token={token}
                  profile={profile?.role === "talent" ? profile : null}
                  onSaved={setProfile}
                />
              ) : (
                <ClientEditor
                  token={token}
                  profile={profile?.role === "client" ? profile : null}
                  onSaved={setProfile}
                />
              )}
            </section>

            <aside className="sheet-side">
              {user.role === "talent" ? (
                <Frames
                  token={token}
                  profile={profile?.role === "talent" ? profile : null}
                  onChanged={load}
                />
              ) : (
                <HiringNote />
              )}
            </aside>
          </div>
        )}
      </main>
    </div>
  );
}

function TalentEditor({
  token,
  profile,
  onSaved,
}: {
  token: string;
  profile: TalentProfile | null;
  onSaved: (profile: Profile) => void;
}) {
  const [form, setForm] = useState({
    display_name: profile?.display_name ?? "",
    city: profile?.city ?? "",
    age: profile?.age ? String(profile.age) : "",
    height_cm: profile?.height_cm != null ? String(profile.height_cm) : "",
    rate_min: profile?.rate_min ? String(profile.rate_min) : "",
    rate_max: profile?.rate_max ? String(profile.rate_max) : "",
    bio: profile?.bio ?? "",
  });
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [saved, setSaved] = useState<string | null>(null);

  function update(key: keyof typeof form) {
    return (event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      setForm((current) => ({ ...current, [key]: event.target.value }));
    };
  }

  async function submit(event: React.FormEvent) {
    event.preventDefault();
    setError(null);
    setSaved(null);

    if (!form.age.trim()) {
      setError("Enter your age. Talent must be 18 or older.");
      return;
    }

    setBusy(true);
    try {
      const updated = await saveTalentProfile(token, {
        talent_type: "model",
        display_name: form.display_name.trim(),
        city: form.city.trim(),
        age: Number(form.age),
        height_cm: form.height_cm.trim() ? Number(form.height_cm) : null,
        bio: form.bio.trim(),
        rate_min: form.rate_min.trim() ? Number(form.rate_min) : 0,
        rate_max: form.rate_max.trim() ? Number(form.rate_max) : 0,
      });
      onSaved(updated);
      setSaved("Sheet saved.");
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setBusy(false);
    }
  }

  return (
    <form className="editor" onSubmit={submit} noValidate>
      <div className="form-grid">
        <Field id="display_name" label="Name on the sheet">
          <input
            id="display_name"
            className="input"
            value={form.display_name}
            onChange={update("display_name")}
            maxLength={80}
            required
          />
        </Field>

        <Field id="city" label="City">
          <input
            id="city"
            className="input"
            value={form.city}
            onChange={update("city")}
            placeholder="Bengaluru"
            required
          />
        </Field>

        <Field id="age" label="Age" hint="18 or older. This is checked when you save.">
          <input
            id="age"
            className="input"
            type="number"
            inputMode="numeric"
            min={18}
            max={100}
            value={form.age}
            onChange={update("age")}
            required
          />
        </Field>

        <Field id="height_cm" label="Height in cm" hint="Optional.">
          <input
            id="height_cm"
            className="input"
            type="number"
            inputMode="numeric"
            min={100}
            max={250}
            value={form.height_cm}
            onChange={update("height_cm")}
          />
        </Field>

        <div className="field span-2">
          <span className="field-label">Day rate in rupees</span>
          <div className="rate-pair">
            <label className="rate-box">
              <span>From</span>
              <input
                className="input"
                type="number"
                inputMode="numeric"
                min={0}
                value={form.rate_min}
                onChange={update("rate_min")}
                aria-label="Minimum day rate in rupees"
              />
            </label>
            <label className="rate-box">
              <span>To</span>
              <input
                className="input"
                type="number"
                inputMode="numeric"
                min={0}
                value={form.rate_max}
                onChange={update("rate_max")}
                aria-label="Maximum day rate in rupees"
              />
            </label>
          </div>
        </div>

        <Field id="bio" label="One line about you" hint={`${form.bio.length}/500`} className="span-2">
          <textarea
            id="bio"
            className="input"
            rows={3}
            maxLength={500}
            value={form.bio}
            onChange={update("bio")}
          />
        </Field>
      </div>

      {error ? (
        <p className="form-error" role="alert">
          {error}
        </p>
      ) : null}

      <div className="form-actions">
        <button className="btn" type="submit" disabled={busy}>
          {busy ? "Saving…" : profile ? "Save sheet" : "Create my sheet"}
        </button>
        {saved ? (
          <span className="save-status" role="status">
            {saved}
          </span>
        ) : null}
      </div>
    </form>
  );
}

function ClientEditor({
  token,
  profile,
  onSaved,
}: {
  token: string;
  profile: ClientProfile | null;
  onSaved: (profile: Profile) => void;
}) {
  const [displayName, setDisplayName] = useState(profile?.display_name ?? "");
  const [orgName, setOrgName] = useState(profile?.org_name ?? "");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [saved, setSaved] = useState<string | null>(null);

  async function submit(event: React.FormEvent) {
    event.preventDefault();
    setError(null);
    setSaved(null);
    setBusy(true);
    try {
      const updated = await saveClientProfile(token, {
        display_name: displayName.trim(),
        org_name: orgName.trim(),
      });
      onSaved(updated);
      setSaved("Sheet saved.");
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setBusy(false);
    }
  }

  return (
    <form className="editor" onSubmit={submit} noValidate>
      <div className="form-grid">
        <Field id="client_name" label="Your name">
          <input
            id="client_name"
            className="input"
            value={displayName}
            onChange={(event) => setDisplayName(event.target.value)}
            maxLength={80}
            required
          />
        </Field>

        <Field id="org_name" label="Studio or brand" hint="Optional, shown on your calls.">
          <input
            id="org_name"
            className="input"
            value={orgName}
            onChange={(event) => setOrgName(event.target.value)}
            maxLength={120}
          />
        </Field>
      </div>

      {error ? (
        <p className="form-error" role="alert">
          {error}
        </p>
      ) : null}

      <div className="form-actions">
        <button className="btn" type="submit" disabled={busy}>
          {busy ? "Saving…" : profile ? "Save sheet" : "Create my sheet"}
        </button>
        {saved ? (
          <span className="save-status" role="status">
            {saved}
          </span>
        ) : null}
      </div>
    </form>
  );
}

function Frames({
  token,
  profile,
  onChanged,
}: {
  token: string;
  profile: TalentProfile | null;
  onChanged: () => Promise<void>;
}) {
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const photos = profile?.photos ?? [];
  const remaining = MAX_PHOTOS - photos.length;

  async function handleFile(event: React.ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    event.target.value = "";
    if (!file) return;

    setError(null);
    setUploading(true);
    try {
      await uploadPhoto(token, file);
      await onChanged();
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setUploading(false);
    }
  }

  return (
    <div className="frames-wrap">
      <div className="frames-head">
        <h2>Frames</h2>
        <span>
          {photos.length} of {MAX_PHOTOS}
        </span>
      </div>

      {!profile ? (
        <p className="side-note">
          Save your details once. Then this side becomes your contact sheet — the first frame is
          what clients see on the deck.
        </p>
      ) : (
        <>
          <div className="frames">
            {photos.map((photo, index) => (
              <figure className="frame" key={photo.id}>
                <img src={assetUrl(photo.url)} alt={`Profile frame ${index + 1}`} loading="lazy" />
                <figcaption>{String(index + 1).padStart(2, "0")}</figcaption>
              </figure>
            ))}

            {remaining > 0 ? (
              <label className={uploading ? "frame frame-add busy" : "frame frame-add"}>
                <span>{uploading ? "Uploading…" : "Add frame"}</span>
                <input
                  className="sr-only"
                  type="file"
                  accept="image/png,image/jpeg,image/webp"
                  onChange={handleFile}
                  disabled={uploading}
                />
              </label>
            ) : null}
          </div>

          <p className="side-note">
            JPG, PNG or WebP. Keep it current — this is what gets shortlisted.
          </p>
        </>
      )}

      {error ? (
        <p className="form-error" role="alert">
          {error}
        </p>
      ) : null}
    </div>
  );
}

function HiringNote() {
  return (
    <div className="hiring-note">
      <h2>Your side of the sheet</h2>
      <p>
        Talent sees who is hiring before they apply. Keep your name and studio accurate so the call
        reads as a real job, not a casting rumour.
      </p>
      <p className="side-note">
        Calls you post will show here with their applicants once jobs go live.
      </p>
    </div>
  );
}
