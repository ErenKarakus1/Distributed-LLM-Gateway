import React from "react";
import { createRoot } from "react-dom/client";
import { Activity, Bot, Clock, DollarSign, RefreshCw, Send, Shield, Zap } from "lucide-react";
import "./styles.css";

type ChatMessage = {
  role: string;
  content: string;
};

type ChatChoice = {
  index: number;
  message: ChatMessage;
  finish_reason: string;
};

type ChatResponse = {
  id: string;
  request_id: string;
  object: string;
  model: string;
  provider: string;
  choices: ChatChoice[];
};

type RequestLog = {
  id: string;
  request_id: string;
  model: string;
  provider?: string;
  status_code: number;
  duration_ms: number;
  prompt_tokens: number;
  completion_tokens: number;
  total_tokens: number;
  input_cost_usd: number;
  output_cost_usd: number;
  total_cost_usd: number;
  created_at: string;
};

type UsageSummary = {
  total_requests: number;
  successful_requests: number;
  failed_requests: number;
  total_tokens: number;
  total_cost_usd: number;
  average_duration_ms: number;
  requests_by_model: Record<string, number>;
  requests_by_provider: Record<string, number>;
};

const models = ["smart", "cheap", "fast", "fail-first", "fail-both", "real"];
const gatewayUrl = import.meta.env.VITE_GATEWAY_URL ?? "/api";

function formatMoney(value: number) {
  return `$${value.toFixed(8)}`;
}

function formatDate(value: string) {
  if (!value) return "-";
  return new Intl.DateTimeFormat(undefined, {
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  }).format(new Date(value));
}

function App() {
  const [apiKey, setApiKey] = React.useState("dev-api-key");
  const [model, setModel] = React.useState("smart");
  const [prompt, setPrompt] = React.useState("Say hello from the dashboard.");
  const [response, setResponse] = React.useState<ChatResponse | null>(null);
  const [logs, setLogs] = React.useState<RequestLog[]>([]);
  const [summary, setSummary] = React.useState<UsageSummary | null>(null);
  const [loading, setLoading] = React.useState(false);
  const [refreshing, setRefreshing] = React.useState(false);
  const [error, setError] = React.useState("");

  const headers = React.useMemo(
    () => ({
      Authorization: `Bearer ${apiKey}`,
      "Content-Type": "application/json",
    }),
    [apiKey],
  );

  async function requestJson<T>(path: string, init?: RequestInit): Promise<T> {
    const res = await fetch(`${gatewayUrl}${path}`, {
      ...init,
      headers: {
        ...headers,
        ...init?.headers,
      },
    });
    const data = await res.json().catch(() => ({}));
    if (!res.ok) {
      throw new Error(data.error ?? `Request failed with status ${res.status}`);
    }
    return data as T;
  }

  async function refreshAdminData(options: { preserveError?: boolean } = {}) {
    setRefreshing(true);
    if (!options.preserveError) {
      setError("");
    }
    try {
      const [nextLogs, nextSummary] = await Promise.all([
        requestJson<RequestLog[]>("/admin/logs", { method: "GET" }),
        requestJson<UsageSummary>("/admin/usage-summary", { method: "GET" }),
      ]);
      setLogs(nextLogs);
      setSummary(nextSummary);
    } catch (err) {
      if (!options.preserveError) {
        setError(err instanceof Error ? err.message : "Failed to refresh dashboard");
      }
    } finally {
      setRefreshing(false);
    }
  }

  async function sendMessage(event: React.FormEvent) {
    event.preventDefault();
    setLoading(true);
    setError("");
    setResponse(null);
    try {
      const result = await requestJson<ChatResponse>("/v1/chat/completions", {
        method: "POST",
        body: JSON.stringify({
          model,
          messages: [{ role: "user", content: prompt }],
        }),
      });
      setResponse(result);
      await refreshAdminData({ preserveError: true });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Chat request failed");
    } finally {
      await refreshAdminData();
      setLoading(false);
    }
  }

  React.useEffect(() => {
    refreshAdminData();
  }, []);

  const successRate = summary && summary.total_requests > 0
    ? (summary.successful_requests / summary.total_requests) * 100
    : 0;

  return (
    <main className="app">
      <header className="topbar">
        <div>
          <h1>Distributed LLM Gateway</h1>
          <p>Chat playground, request logs, and usage telemetry for the gateway.</p>
        </div>
        <div className="api-key">
          <label htmlFor="api-key">API key</label>
          <input
            id="api-key"
            value={apiKey}
            onChange={(event) => setApiKey(event.target.value)}
            spellCheck={false}
          />
        </div>
      </header>

      {error && <div className="notice">{error}</div>}

      <section className="layout">
        <section className="panel playground">
          <div className="section-title">
            <Bot size={18} />
            <h2>Chat Playground</h2>
          </div>

          <form onSubmit={sendMessage} className="chat-form">
            <label htmlFor="model">Model</label>
            <select id="model" value={model} onChange={(event) => setModel(event.target.value)}>
              {models.map((item) => (
                <option key={item} value={item}>
                  {item}
                </option>
              ))}
            </select>

            <label htmlFor="prompt">Message</label>
            <textarea
              id="prompt"
              value={prompt}
              onChange={(event) => setPrompt(event.target.value)}
              rows={8}
            />

            <button className="primary-button" type="submit" disabled={loading || !prompt.trim()}>
              <Send size={17} />
              {loading ? "Sending" : "Send"}
            </button>
          </form>

          <div className="response-box">
            <div className="response-meta">
              <span>{response?.provider ?? "No provider yet"}</span>
              <span>{response?.request_id ?? "No request yet"}</span>
            </div>
            <p>{response?.choices[0]?.message.content ?? "Send a request to see the model response here."}</p>
          </div>
        </section>

        <section className="dashboard">
          <div className="summary-grid">
            <Metric icon={<Activity size={18} />} label="Requests" value={summary?.total_requests ?? 0} />
            <Metric icon={<Shield size={18} />} label="Success Rate" value={`${successRate.toFixed(1)}%`} />
            <Metric icon={<Zap size={18} />} label="Tokens" value={summary?.total_tokens ?? 0} />
            <Metric icon={<DollarSign size={18} />} label="Cost" value={formatMoney(summary?.total_cost_usd ?? 0)} />
            <Metric icon={<Clock size={18} />} label="Avg Latency" value={`${(summary?.average_duration_ms ?? 0).toFixed(1)} ms`} />
          </div>

          <section className="panel split-panel">
            <div>
              <div className="section-title">
                <Activity size={18} />
                <h2>Models</h2>
              </div>
              <Breakdown data={summary?.requests_by_model ?? {}} />
            </div>
            <div>
              <div className="section-title">
                <Zap size={18} />
                <h2>Providers</h2>
              </div>
              <Breakdown data={summary?.requests_by_provider ?? {}} />
            </div>
          </section>

          <section className="panel logs-panel">
            <div className="logs-header">
              <div className="section-title">
                <Clock size={18} />
                <h2>Recent Requests</h2>
              </div>
              <button className="icon-button" type="button" onClick={() => refreshAdminData()} disabled={refreshing} title="Refresh">
                <RefreshCw size={17} />
              </button>
            </div>

            <div className="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th>Time</th>
                    <th>Model</th>
                    <th>Provider</th>
                    <th>Status</th>
                    <th>Latency</th>
                    <th>Tokens</th>
                    <th>Cost</th>
                  </tr>
                </thead>
                <tbody>
                  {logs.map((log) => (
                    <tr key={log.id}>
                      <td>{formatDate(log.created_at)}</td>
                      <td>{log.model}</td>
                      <td>{log.provider ?? "-"}</td>
                      <td>
                        <span className={log.status_code >= 200 && log.status_code < 300 ? "status-ok" : "status-error"}>
                          {log.status_code}
                        </span>
                      </td>
                      <td>{log.duration_ms} ms</td>
                      <td>{log.total_tokens}</td>
                      <td>{formatMoney(log.total_cost_usd)}</td>
                    </tr>
                  ))}
                  {logs.length === 0 && (
                    <tr>
                      <td colSpan={7} className="empty-row">
                        No request logs yet.
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </section>
        </section>
      </section>
    </main>
  );
}

function Metric({ icon, label, value }: { icon: React.ReactNode; label: string; value: React.ReactNode }) {
  return (
    <div className="metric">
      <div className="metric-icon">{icon}</div>
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

function Breakdown({ data }: { data: Record<string, number> }) {
  const entries = Object.entries(data);
  const max = Math.max(...entries.map(([, value]) => value), 1);

  if (entries.length === 0) {
    return <p className="muted">No data yet.</p>;
  }

  return (
    <div className="breakdown">
      {entries.map(([label, value]) => (
        <div className="bar-row" key={label}>
          <div className="bar-label">
            <span>{label}</span>
            <strong>{value}</strong>
          </div>
          <div className="bar-track">
            <div className="bar-fill" style={{ width: `${(value / max) * 100}%` }} />
          </div>
        </div>
      ))}
    </div>
  );
}

createRoot(document.getElementById("root")!).render(<App />);
