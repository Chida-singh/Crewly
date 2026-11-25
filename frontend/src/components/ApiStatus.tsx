import { useEffect, useState } from "react";

import { getHealth, type Health } from "../lib/api";

export function ApiStatus() {
  const [health, setHealth] = useState<Health | null>(null);

  useEffect(() => {
    let live = true;
    void getHealth().then((result) => {
      if (live) setHealth(result);
    });
    return () => {
      live = false;
    };
  }, []);

  const label = !health
    ? "Checking API"
    : health.status === "offline"
      ? "API offline"
      : health.database === "down"
        ? "Database down"
        : "API up";

  const tone = !health ? "idle" : label === "API up" ? "ok" : "warn";

  return (
    <span className={`api-status ${tone}`}>
      <span className="dot" aria-hidden="true" />
      {label}
    </span>
  );
}
