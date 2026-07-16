import { useState } from 'react'

export function SettingsPage() {
  const [provider, setProvider] = useState(() => localStorage.getItem('ai_provider') ?? 'openai')
  const [apiKey, setApiKey] = useState(() => localStorage.getItem('ai_api_key') ?? '')
  const [saved, setSaved] = useState(false)

  function save() {
    localStorage.setItem('ai_provider', provider)
    localStorage.setItem('ai_api_key', apiKey)
    setSaved(true)
    setTimeout(() => setSaved(false), 2000)
  }

  return (
    <div className="page-narrow">
      <h2 style={{ fontSize: '1.5rem', fontWeight: 800, marginBottom: 8 }}>AI Settings</h2>
      <p className="text-sm text-muted" style={{ marginBottom: 28 }}>
        Your API key is stored only in this browser — never on our servers, except during quiz generation.
      </p>

      <div className="card">
        <div className="form-stack">
          <div className="form-group">
            <label className="form-label">AI Provider</label>
            <select
              className="form-input"
              value={provider}
              onChange={(e) => setProvider(e.target.value)}
            >
              <option value="openai">OpenAI — GPT-4o mini</option>
              <option value="gemini">Google — Gemini 2.0 Flash</option>
              <option value="claude">Anthropic — Claude Haiku</option>
            </select>
          </div>

          <div className="form-group">
            <label className="form-label">API Key</label>
            <input
              className="form-input"
              type="password"
              placeholder="Paste your API key..."
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
            />
          </div>

          {saved && <div className="alert alert-success">Settings saved!</div>}

          <button className="btn btn-primary btn-full" onClick={save}>
            Save Settings
          </button>
        </div>
      </div>

      <div className="card" style={{ marginTop: 16, padding: 16 }}>
        <p className="text-sm text-muted" style={{ lineHeight: 1.7 }}>
          <strong>OpenAI:</strong> Get key at platform.openai.com<br />
          <strong>Gemini:</strong> Get key at aistudio.google.com<br />
          <strong>Claude:</strong> Get key at console.anthropic.com
        </p>
      </div>
    </div>
  )
}
