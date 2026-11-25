import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

import { errorMessage, requestOtp, verifyOtp, type Role } from "../lib/api";
import { useSession } from "../lib/session";

const RESEND_SECONDS = 30;

export function AuthTicket() {
  const navigate = useNavigate();
  const { login } = useSession();

  const [role, setRole] = useState<Role>("talent");
  const [step, setStep] = useState<"phone" | "code">("phone");
  const [phone, setPhone] = useState("");
  const [code, setCode] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [notice, setNotice] = useState<string | null>(null);
  const [cooldown, setCooldown] = useState(0);

  useEffect(() => {
    if (cooldown <= 0) return;
    const timer = window.setTimeout(() => setCooldown((value) => value - 1), 1000);
    return () => window.clearTimeout(timer);
  }, [cooldown]);

  async function sendCode(event?: React.FormEvent) {
    event?.preventDefault();
    setError(null);
    setBusy(true);
    try {
      await requestOtp(phone);
      setStep("code");
      setNotice("Code sent. In local development it is printed in the API log.");
      setCooldown(RESEND_SECONDS);
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setBusy(false);
    }
  }

  async function verify(event: React.FormEvent) {
    event.preventDefault();
    setError(null);
    setBusy(true);
    try {
      const result = await verifyOtp(phone, code.trim(), role);
      login({ token: result.token, user: result.user });
      navigate("/sheet", { replace: true });
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setBusy(false);
    }
  }

  function resetPhone() {
    setStep("phone");
    setCode("");
    setNotice(null);
    setError(null);
  }

  return (
    <div className="ticket">
      <div className="ticket-stub">
        <span>Check in</span>
        <span>{step === "phone" ? "Step 1 of 2" : "Step 2 of 2"}</span>
      </div>

      <div className="ticket-body">
        <div className="role-switch" role="group" aria-label="Account type">
          <button
            type="button"
            aria-pressed={role === "talent"}
            onClick={() => setRole("talent")}
            disabled={step === "code"}
          >
            Talent
          </button>
          <button
            type="button"
            aria-pressed={role === "client"}
            onClick={() => setRole("client")}
            disabled={step === "code"}
          >
            Client
          </button>
        </div>

        <p className="ticket-hint">
          {step === "phone"
            ? "New numbers pick a side. Existing accounts keep the role they signed up with."
            : "Enter the code we sent to your phone."}
        </p>

        <form className="ticket-form" onSubmit={step === "phone" ? sendCode : verify}>
          {step === "phone" ? (
            <div className="field">
              <label htmlFor="phone">Phone number</label>
              <div className="phone-input">
                <span aria-hidden="true">+91</span>
                <input
                  id="phone"
                  name="phone"
                  type="tel"
                  inputMode="numeric"
                  autoComplete="tel-national"
                  maxLength={10}
                  placeholder="98765 43210"
                  value={phone}
                  onChange={(event) => setPhone(event.target.value.replace(/\D/g, "").slice(0, 10))}
                  autoFocus
                />
              </div>
            </div>
          ) : (
            <>
              <p className="ticket-phone">
                +91 {phone}
                <button type="button" className="link" onClick={resetPhone}>
                  Change
                </button>
              </p>
              <div className="field">
                <label htmlFor="code">6-digit code</label>
                <input
                  id="code"
                  name="code"
                  className="code-input"
                  inputMode="numeric"
                  autoComplete="one-time-code"
                  maxLength={6}
                  placeholder="000000"
                  value={code}
                  onChange={(event) => setCode(event.target.value.replace(/\D/g, "").slice(0, 6))}
                  autoFocus
                />
              </div>
            </>
          )}

          {error ? (
            <p className="form-error" role="alert">
              {error}
            </p>
          ) : null}
          {notice && !error ? <p className="form-notice">{notice}</p> : null}

          <button
            className="btn"
            type="submit"
            disabled={busy || (step === "phone" ? phone.length < 10 : code.length < 6)}
          >
            {busy ? "Working…" : step === "phone" ? "Send code" : "Verify and continue"}
          </button>

          {step === "code" ? (
            <button
              type="button"
              className="btn ghost small"
              disabled={busy || cooldown > 0}
              onClick={() => void sendCode()}
            >
              {cooldown > 0 ? `Resend in ${cooldown}s` : "Resend code"}
            </button>
          ) : null}
        </form>

        <p className="ticket-fine">Phone verified. 18+ only. One number, one account.</p>
      </div>
    </div>
  );
}
