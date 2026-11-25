import { Link, Navigate } from "react-router-dom";

import { ApiStatus } from "../components/ApiStatus";
import { AuthTicket } from "../components/AuthTicket";
import { useSession } from "../lib/session";

const steps = [
  {
    number: "01",
    name: "Talent",
    copy: "Build a sheet with photos, your city and your day rate. Apply to the calls you want.",
  },
  {
    number: "02",
    name: "Clients",
    copy: "Post a real brief. Shortlist from the people who already said yes to the work.",
  },
  {
    number: "03",
    name: "Match",
    copy: "Chat opens only when both sides agree. No cold calls, no manager's cut.",
  },
];

export function LandingPage() {
  const { session } = useSession();
  if (session) return <Navigate to="/sheet" replace />;

  return (
    <div className="page">
      <header className="topbar">
        <Link to="/" className="wordmark" aria-label="Crewly home">
          CREWLY<span aria-hidden="true">.</span>
        </Link>
        <ApiStatus />
      </header>

      <main>
        <div className="hero">
          <section className="poster" aria-labelledby="poster-title">
            <h1 id="poster-title">
              <span className="outline">Open</span>
              <br />
              Call
            </h1>

            <div className="poster-foot">
              <p className="poster-copy">
                Crewly is the casting sheet for India&rsquo;s creative workforce. Build your sheet,
                get shortlisted, and keep your day rate.
              </p>

              <div className="poster-tags">
                <span>Models</span>
                <span>Bengaluru</span>
                <span>18+ only</span>
              </div>
            </div>

            <div className="stamp" aria-hidden="true">
              <span>Open call</span>
              <span>Bengaluru</span>
              <span>Est. 2026</span>
            </div>
          </section>

          <section className="checkin" aria-label="Sign in">
            <AuthTicket />
          </section>
        </div>

        <section className="process" aria-labelledby="process-title">
          <h2 id="process-title" className="process-title">
            How the sheet works
          </h2>
          {steps.map((step) => (
            <div className="process-row" key={step.number}>
              <span className="process-num">{step.number}</span>
              <span className="process-name">{step.name}</span>
              <span className="process-copy">{step.copy}</span>
            </div>
          ))}
        </section>
      </main>

      <footer className="footer">
        <span>Crewly — direct hiring for creative work.</span>
        <span>Bengaluru · 18+ · Phone verified</span>
      </footer>
    </div>
  );
}
