<script>
  import { onMount } from 'svelte';
  
  // Wails-generated JS bindings
  import { 
    AddDownload, 
    GetJobs, 
    CancelJob, 
    SubmitSteamCode, 
    CancelSteamCodePrompt, 
    GetConfig, 
    SaveConfig,
    LoginSteam,
    CheckSteamCmd,
    ForceInstallSteamCmd,
    RetryJob,
    DeleteJob
  } from '../wailsjs/go/main/App';
  import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime';

  // Component State
  let config = { steamCmdPath: '', outputDir: '', autoUpdate: true, username: 'anonymous' };
  let workshopUrl = '';
  let manualAppId = '';
  let jobs = [];
  
  // Settings Modal State
  let showConfigModal = false;
  let steamPassword = '';
  let isSavingSettings = false;
  let isLoggingIn = false;
  let loginStatusMsg = '';
  let isCheckingSteamCmd = false;

  // Active Steam Guard input modal details
  let showAuthModal = false;
  let authPromptType = ''; // 'email' or '2fa'
  let steamCode = '';

  // Stream logs
  let activeJobLogs = [];
  let logContainer;

  onMount(async () => {
    // 1. Fetch config and active jobs list
    await reloadConfig();
    await fetchJobs();

    // 2. Register Wails events
    EventsOn('job:status', (updatedJob) => {
      if (!jobs) jobs = [];
      const idx = jobs.findIndex(j => j.id === updatedJob.id);
      if (idx !== -1) {
        jobs[idx] = updatedJob;
      } else {
        jobs = [updatedJob, ...jobs];
      }
      fetchJobs();
    });

    EventsOn('job:progress', (event) => {
      if (!jobs) jobs = [];
      const idx = jobs.findIndex(j => j.id === event.jobId);
      if (idx !== -1) {
        jobs[idx].progress = event.progress;
      }
    });

    EventsOn('steamcmd:log', (logEntry) => {
      activeJobLogs = [...activeJobLogs, {
        id: Math.random().toString(),
        time: logEntry.timestamp.slice(11, 19),
        text: logEntry.message,
        isError: logEntry.isError
      }];
      if (activeJobLogs.length > 500) {
        activeJobLogs.shift();
      }
      
      // Auto-scroll to bottom of logs
      setTimeout(() => {
        if (logContainer) {
          logContainer.scrollTop = logContainer.scrollHeight;
        }
      }, 30);
    });

    EventsOn('steamcmd:auth_required', (type) => {
      authPromptType = type;
      showAuthModal = true;
      steamCode = '';
    });

    EventsOn('job:deleted', (jobId) => {
      jobs = (jobs || []).filter(j => j.id !== jobId);
    });

    return () => {
      EventsOff('job:status');
      EventsOff('job:progress');
      EventsOff('steamcmd:log');
      EventsOff('steamcmd:auth_required');
      EventsOff('job:deleted');
    };
  });

  async function reloadConfig() {
    config = await GetConfig();
  }

  async function updateConfig() {
    isSavingSettings = true;
    try {
      await SaveConfig(config);
      showConfigModal = false;
    } catch (err) {
      alert(`Error saving config: ${err}`);
    } finally {
      isSavingSettings = false;
    }
  }

  async function triggerSteamLogin() {
    if (!config.username || config.username === 'anonymous') {
      alert('Please enter a valid Steam Username (different from anonymous) to authenticate.');
      return;
    }
    isLoggingIn = true;
    loginStatusMsg = 'Initiating SteamCMD login... Check logs below.';
    try {
      await LoginSteam(config.username, steamPassword);
      loginStatusMsg = '';
      steamPassword = '';
      alert('Steam authentication completed successfully! Credentials cached in SteamCMD.');
    } catch (err) {
      loginStatusMsg = '';
      alert(`Steam login failed: ${err}`);
    } finally {
      isLoggingIn = false;
      await reloadConfig();
    }
  }

  async function handleVerifySteamCmd() {
    isCheckingSteamCmd = true;
    try {
      const status = await CheckSteamCmd();
      if (status === 'found') {
        alert('SteamCMD executable found and verified!');
      } else {
        const fetchNow = confirm('SteamCMD was not found at this location. Would you like to download and install it now?');
        if (fetchNow) {
          isLoggingIn = true;
          loginStatusMsg = 'Downloading and installing SteamCMD Client... Please wait.';
          try {
            await ForceInstallSteamCmd();
            alert('SteamCMD installed successfully!');
          } catch (installErr) {
            alert(`Installation failed: ${installErr}`);
          } finally {
            isLoggingIn = false;
          }
        }
      }
    } catch (err) {
      alert(`Error verifying path: ${err}`);
    } finally {
      isCheckingSteamCmd = false;
    }
  }

  async function fetchJobs() {
    try {
      const result = await GetJobs();
      jobs = result || [];
    } catch (err) {
      jobs = [];
    }
  }

  async function handleSubmitDownload() {
    try {
      await AddDownload(workshopUrl, manualAppId);
      workshopUrl = '';
      manualAppId = '';
      await fetchJobs();
    } catch (err) {
      alert(`Error submitting job: ${err}`);
    }
  }

  async function handleCancel(jobId) {
    await CancelJob(jobId);
    await fetchJobs();
  }

  async function handleRetry(jobId) {
    try {
      await RetryJob(jobId);
      await fetchJobs();
    } catch (err) {
      alert(`Retry failed: ${err}`);
    }
  }

  async function handleDelete(jobId) {
    await DeleteJob(jobId);
    await fetchJobs();
  }

  async function submitCode() {
    await SubmitSteamCode(steamCode);
    showAuthModal = false;
  }

  async function cancelCode() {
    await CancelSteamCodePrompt();
    showAuthModal = false;
  }
</script>

<main class="app-workspace">
  <!-- Top Bar Decorator Pattern -->
  <header class="top-nav">
    <div class="brand">
      <span class="pulse-icon"></span>
      <span class="brand-title">SteamCMD Workshop Downloader</span>
    </div>
    <div class="navbar-actions">
      <button class="btn-secondary" on:click={() => showConfigModal = true}>
        ⚙️ Settings & Auth
      </button>
    </div>
  </header>

  <div class="workspace-body">
    <!-- Left Column: Input Forms & Stats -->
    <aside class="left-pane">
      <section class="card">
        <div class="card-header">
          <h3>New Download Task</h3>
        </div>
        <div class="card-body">
          <form on:submit|preventDefault={handleSubmitDownload} class="download-form">
            <div class="form-group">
              <label for="urlInput">Workshop Item URL or ID</label>
              <input 
                id="urlInput" 
                type="text" 
                placeholder="https://steamcommunity.com/sharedfiles/..." 
                bind:value={workshopUrl} 
                required 
              />
            </div>
            
            <div class="form-group">
              <label for="appIdInput">Manual AppID (if query appid failed)</label>
              <input 
                id="appIdInput" 
                type="text" 
                placeholder="e.g. 281990" 
                bind:value={manualAppId} 
              />
            </div>

            <button type="submit" class="btn-primary btn-block">Add to Queue</button>
          </form>
        </div>
      </section>

      <section class="card stats-card">
        <div class="card-body">
          <div class="status-summary">
            <div class="stat-item">
              <span class="stat-val">{(jobs || []).length}</span>
              <span class="stat-lbl">Total Tasks</span>
            </div>
            <div class="stat-item">
              <span class="stat-val">{(jobs || []).filter(j => j && j.status === 'running').length}</span>
              <span class="stat-lbl">Active</span>
            </div>
            <div class="stat-item">
              <span class="stat-val">{(jobs || []).filter(j => j && j.status === 'success').length}</span>
              <span class="stat-lbl">Success</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Connection Status Overview -->
      <section class="card details-card">
        <div class="card-header">
          <h3>Session Account</h3>
        </div>
        <div class="card-body">
          <div class="session-info">
            <div class="session-row">
              <span class="lbl font-semibold">Active User:</span>
              <span class="val highlight-text">{config.username}</span>
            </div>
            <div class="session-row">
              <span class="lbl">Engine Status:</span>
              <span class="val text-green">Ready</span>
            </div>
          </div>
        </div>
      </section>
    </aside>

    <!-- Right Column: Queue (Scrollable) & Terminal Logs (Scrollable) -->
    <section class="right-pane">
      <!-- Queue Panel -->
      <div class="card queue-card">
        <div class="card-header flex-header">
          <h3>Download Queue</h3>
          <span class="badge">{(jobs || []).length} items</span>
        </div>
        <div class="card-body scrollable-content">
          {#if !(jobs && jobs.length)}
            <div class="empty-state">
              <p>Download queue is empty</p>
              <small>Submit a Steam Workshop URL above to begin downloading.</small>
            </div>
          {:else}
            <div class="queue-list">
              {#each jobs || [] as job (job.id)}
                <div class="job-item status-{job.status}">
                  <div class="job-card-top">
                    {#if job.previewUrl}
                      <img src={job.previewUrl} alt={job.title} class="job-preview-img" />
                    {/if}
                    <div class="job-details">
                      <div class="job-title-row">
                        <strong class="job-title-text" title={job.title || `Workshop ID: ${job.workshopId}`}>
                          {job.title || `Workshop ID: ${job.workshopId}`}
                        </strong>
                        <span class="job-status-badge">{job.status}</span>
                      </div>
                      <div class="job-ids">
                        <span class="text-secondary">AppID: {job.appId} | ID: {job.workshopId}</span>
                      </div>
                    </div>
                  </div>

                  <div class="progress-section">
                    <div class="progress-bar-bg">
                      <div class="progress-bar-fill" style="width: {job.progress}%"></div>
                    </div>
                    <span class="progress-percentage">{job.progress.toFixed(1)}%</span>
                  </div>

                  {#if job.errorMsg}
                    <div class="job-error-msg">⚠️ {job.errorMsg}</div>
                  {/if}

                  <div class="job-actions">
                    {#if job.status === 'queued' || job.status === 'running'}
                      <button class="btn-danger-text" on:click={() => handleCancel(job.id)}>Cancel</button>
                    {:else}
                      {#if job.status === 'failed'}
                        <button class="btn-retry-text" on:click={() => handleRetry(job.id)}>🔄 Retry</button>
                      {/if}
                      <button class="btn-delete-text" on:click={() => handleDelete(job.id)}>🗑️ Remove</button>
                    {/if}
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      </div>

      <!-- Live Logs Console -->
      <div class="card console-card">
        <div class="card-header console-header">
          <h3>SteamCMD Live Logs</h3>
          <button class="btn-icon-text" on:click={() => activeJobLogs = []}>🗑️ Clear Logs</button>
        </div>
        <div class="console-body" bind:this={logContainer}>
          {#if activeJobLogs.length === 0}
            <div class="console-placeholder">Listening for SteamCMD console events...</div>
          {:else}
            {#each activeJobLogs as log (log.id)}
              <div class="line" class:error-line={log.isError}>
                <span class="timestamp">[{log.time}]</span>
                <span class="msg">{log.text}</span>
              </div>
            {/each}
          {/if}
        </div>
      </div>
    </section>
  </div>

  <!-- Settings & Authentication Modal Overlay -->
  {#if showConfigModal}
    <div class="modal-backdrop">
      <div class="modal-panel settings-modal">
        <div class="modal-header">
          <h3>Settings & Steam Auth</h3>
        </div>
        <div class="modal-body modal-scroll">
          <div class="modal-section-title">📂 Paths & Configurations</div>
          
          <div class="form-group">
            <label for="steamPath">SteamCMD Path</label>
            <div class="input-with-button">
              <input id="steamPath" type="text" bind:value={config.steamCmdPath} placeholder="C:\path\to\steamcmd.exe" />
              <button class="btn-check" type="button" on:click={handleVerifySteamCmd} disabled={isCheckingSteamCmd}>
                {isCheckingSteamCmd ? 'Verifying...' : '🔍 Check'}
              </button>
            </div>
          </div>
          
          <div class="form-group">
            <label for="outDir">Downloads Directory</label>
            <input id="outDir" type="text" bind:value={config.outputDir} placeholder="C:\Downloads" />
          </div>

          <div class="modal-section-title">🔑 Steam Authentication</div>
          <p class="description-text">
            Logging in caches a persistence token (Sentry) locally inside your SteamCMD environment. 
            Once authenticated, future downloads using this account will not prompt for a password.
          </p>

          <div class="form-group">
            <label for="username">Steam Account Username</label>
            <input id="username" type="text" bind:value={config.username} placeholder="anonymous" disabled={isLoggingIn} />
            <small class="hint">Use "anonymous" to download open-access items without login.</small>
          </div>

          {#if config.username !== 'anonymous'}
            <div class="form-group">
              <label for="password">Steam Account Password</label>
              <input id="password" type="password" bind:value={steamPassword} placeholder="••••••••" disabled={isLoggingIn} />
              <small class="hint">Password is never saved. It is only sent to SteamCMD for the initial login handshake.</small>
            </div>
          {/if}

          {#if isLoggingIn}
            <div class="login-spinner-container">
              <span class="spinner"></span>
              <p class="spinner-msg">{loginStatusMsg}</p>
            </div>
          {:else if config.username !== 'anonymous'}
            <button class="btn-auth btn-block" on:click={triggerSteamLogin}>
              🚀 Connect & Authenticate Account
            </button>
          {/if}
        </div>
        <div class="modal-footer border-top">
          <button class="btn-primary" on:click={updateConfig} disabled={isLoggingIn || isSavingSettings}>
            {isSavingSettings ? 'Saving...' : 'Save Paths'}
          </button>
          <button class="btn-secondary" on:click={() => { showConfigModal = false; steamPassword = ''; }} disabled={isLoggingIn}>
            Cancel
          </button>
        </div>
      </div>
    </div>
  {/if}

  <!-- Modal Dialog: Steam Guard Authenticator -->
  {#if showAuthModal}
    <div class="modal-backdrop z-high">
      <div class="modal-panel auth-code-modal">
        <div class="modal-header">
          <h3>Steam Guard Verification</h3>
        </div>
        <div class="modal-body">
          <p>
            Your account requires approval. Please supply the Steam Guard code received via 
            <strong>{authPromptType === 'email' ? 'Email Address' : 'Mobile App (2FA)'}</strong>.
          </p>
          <input 
            type="text" 
            bind:value={steamCode} 
            placeholder="e.g. A1B2C" 
            class="input-code" 
            autofocus 
          />
        </div>
        <div class="modal-footer">
          <button class="btn-primary" on:click={submitCode}>Submit Code</button>
          <button class="btn-secondary" on:click={cancelCode}>Cancel</button>
        </div>
      </div>
    </div>
  {/if}
</main>

<style>
  :global(html), :global(body) {
    margin: 0;
    padding: 0;
    width: 100%;
    height: 100%;
    overflow: hidden;
    background-color: #0b0f19;
    color: #e2e8f0;
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  }

  .input-with-button {
    display: flex;
    gap: 8px;
    width: 100%;
  }

  .btn-check {
    background: #334155;
    color: #f8fafc;
    border: 1px solid #475569;
    white-space: nowrap;
    padding: 10px 16px;
  }

  .btn-check:hover:not(:disabled) {
    background: #475569;
  }

  label {
    font-weight: 600;
    font-size: 12px;
    color: #94a3b8;
    margin-bottom: 2px;
  }

  input[type="text"], input[type="password"] {
    background: #0f172a;
    border: 1px solid #334155;
    border-radius: 6px;
    padding: 10px 12px;
    color: #f8fafc;
    font-size: 14px;
    width: 100%;
    box-sizing: border-box;
    transition: all 0.2s;
  }

  input[type="text"]:focus, input[type="password"]:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
  }

  input:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  button {
    font-family: inherit;
    font-size: 13px;
    font-weight: 600;
    border-radius: 6px;
    padding: 8px 16px;
    cursor: pointer;
    transition: all 0.2s;
    border: none;
  }

  button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .btn-primary {
    background: #2563eb;
    color: #ffffff;
  }

  .btn-primary:hover:not(:disabled) {
    background: #1d4ed8;
  }

  .btn-primary.btn-block {
    width: 100%;
    display: block;
    padding: 11px;
  }

  .btn-secondary {
    background: #1e293b;
    border: 1px solid #334155;
    color: #cbd5e1;
  }

  .btn-secondary:hover:not(:disabled) {
    background: #334155;
    color: #f8fafc;
  }

  .btn-auth {
    background: #10b981;
    color: #ffffff;
    padding: 11px;
    margin-top: 10px;
  }

  .btn-auth:hover:not(:disabled) {
    background: #059669;
  }

  .btn-danger-text {
    background: transparent;
    color: #f87171;
    font-size: 12px;
    padding: 4px 8px;
  }

  .btn-danger-text:hover {
    background: rgba(239, 68, 68, 0.15);
    border-radius: 4px;
  }

  .btn-icon-text {
    background: transparent;
    color: #94a3b8;
    font-size: 12px;
    padding: 4px 8px;
  }

  .btn-icon-text:hover {
    color: #f8fafc;
    background: rgba(255, 255, 255, 0.05);
  }

  /* Fullpage Layout Core */
  .app-workspace {
    display: flex;
    flex-direction: column;
    width: 100vw;
    height: 100vh;
    box-sizing: border-box;
  }

  .top-nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: #111827;
    border-bottom: 1px solid #1f2937;
    padding: 12px 20px;
    height: 54px;
    box-sizing: border-box;
  }

  .brand {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .pulse-icon {
    width: 8px;
    height: 8px;
    background: #10b981;
    border-radius: 50%;
    box-shadow: 0 0 8px #10b981;
  }

  .brand-title {
    font-size: 15px;
    font-weight: 700;
    letter-spacing: 0.5px;
    color: #f8fafc;
  }

  .workspace-body {
    display: flex;
    flex: 1;
    height: calc(100vh - 54px);
    overflow: hidden;
  }

  .left-pane {
    width: 320px;
    background: #0b0f19;
    border-right: 1px solid #1f2937;
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 20px;
    box-sizing: border-box;
    overflow-y: auto;
  }

  .right-pane {
    flex: 1;
    background: #090d16;
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 20px;
    box-sizing: border-box;
    overflow: hidden;
  }

  /* Cards Design System */
  .card {
    background: #111827;
    border: 1px solid #1f2937;
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .card-header {
    background: #151e2e;
    padding: 12px 16px;
    border-bottom: 1px solid #1f2937;
  }

  .card-header h3 {
    margin: 0;
    font-size: 13px;
    font-weight: 700;
    color: #f1f5f9;
  }

  .card-body {
    padding: 16px;
  }

  .flex-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .badge {
    background: #1e293b;
    color: #94a3b8;
    padding: 2px 8px;
    border-radius: 12px;
    font-size: 11px;
    font-weight: 600;
  }

  /* Form layouts */
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-bottom: 14px;
  }

  .hint {
    font-size: 11px;
    color: #64748b;
  }

  /* Stats Row */
  .status-summary {
    display: flex;
    justify-content: space-between;
    text-align: center;
  }

  .stat-item {
    display: flex;
    flex-direction: column;
  }

  .stat-val {
    font-size: 20px;
    font-weight: 800;
    color: #e2e8f0;
  }

  .stat-lbl {
    font-size: 11px;
    color: #64748b;
    margin-top: 2px;
  }

  .session-info {
    display: flex;
    flex-direction: column;
    gap: 8px;
    font-size: 13px;
  }

  .session-row {
    display: flex;
    justify-content: space-between;
  }

  .highlight-text {
    color: #3b82f6;
    font-weight: 700;
  }

  .text-green {
    color: #10b981;
    font-weight: 700;
  }

  /* Scrollable Queue Layout */
  .queue-card {
    flex: 2;
    min-height: 200px;
    overflow: hidden;
  }

  .scrollable-content {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
  }

  .empty-state {
    text-align: center;
    color: #64748b;
    padding: 40px 0;
  }

  .empty-state p {
    margin: 0 0 6px 0;
    font-weight: 600;
    font-size: 14px;
  }

  .queue-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .job-item {
    background: #1e293b;
    border: 1px solid #334155;
    border-left: 4px solid #94a3b8;
    border-radius: 6px;
    padding: 12px 16px;
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .job-item.status-queued {
    border-left-color: #fbbf24;
  }

  .job-item.status-running {
    border-left-color: #3b82f6;
  }

  .job-item.status-success {
    border-left-color: #10b981;
  }

  .job-item.status-failed {
    border-left-color: #ef4444;
  }

  .job-meta {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
  }

  .job-ids {
    display: flex;
    flex-direction: column;
    gap: 2px;
    font-size: 13px;
  }

  .job-status-badge {
    font-size: 10px;
    text-transform: uppercase;
    font-weight: 700;
    background: rgba(255, 255, 255, 0.05);
    padding: 2px 6px;
    border-radius: 4px;
  }

  .progress-section {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .progress-bar-bg {
    flex: 1;
    background: #0f172a;
    height: 8px;
    border-radius: 4px;
    overflow: hidden;
  }

  .progress-bar-fill {
    background: #3b82f6;
    height: 100%;
    border-radius: 4px;
    transition: width 0.3s ease;
  }

  .job-item.status-success .progress-bar-fill {
    background: #10b981;
  }

  .progress-percentage {
    font-size: 12px;
    font-weight: 700;
    color: #f1f5f9;
    min-width: 36px;
    text-align: right;
  }

  .job-error-msg {
    background: rgba(239, 68, 68, 0.1);
    color: #fca5a5;
    padding: 6px 10px;
    border-radius: 4px;
    font-size: 12px;
  }

  .job-actions {
    display: flex;
    justify-content: flex-end;
  }

  /* Console Layout */
  .console-card {
    height: 240px;
    min-height: 240px;
    background: #090d16;
  }

  .console-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: #0f1624;
    padding: 8px 16px;
  }

  .console-body {
    flex: 1;
    background: #020617;
    padding: 12px;
    overflow-y: auto;
    font-family: 'Fira Code', 'Cascadia Code', monospace;
    font-size: 11px;
    line-height: 1.6;
    display: flex;
    flex-direction: column;
    gap: 4px;
    color: #a7f3d0;
    border-top: 1px solid #1f2937;
  }

  .console-placeholder {
    color: #475569;
    text-align: center;
    margin: auto;
    font-family: sans-serif;
  }

  .line {
    word-break: break-all;
    white-space: pre-wrap;
  }

  .line.error-line {
    color: #fca5a5;
  }

  .timestamp {
    color: #4b5563;
    margin-right: 6px;
  }

  /* Modals Overlay */
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(4, 6, 15, 0.88);
    backdrop-filter: blur(5px);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .z-high {
    z-index: 1100;
  }

  .modal-panel {
    background: #111827;
    border: 1px solid #1f2937;
    border-radius: 8px;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.6);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .settings-modal {
    width: 520px;
    max-height: 85vh;
  }

  .auth-code-modal {
    width: 380px;
    border-color: #3b82f6;
  }

  .modal-header h3 {
    margin: 0;
    font-size: 15px;
    font-weight: 700;
    color: #f8fafc;
  }

  .modal-body {
    padding: 20px;
  }

  .modal-scroll {
    overflow-y: auto;
    max-height: calc(85vh - 120px);
  }

  .modal-section-title {
    font-size: 11px;
    font-weight: 800;
    text-transform: uppercase;
    color: #3b82f6;
    letter-spacing: 0.5px;
    margin: 10px 0 14px 0;
    border-bottom: 1px solid #1f2937;
    padding-bottom: 6px;
  }

  .description-text {
    font-size: 12px;
    line-height: 1.5;
    color: #94a3b8;
    margin-bottom: 16px;
  }

  .login-spinner-container {
    display: flex;
    align-items: center;
    gap: 12px;
    background: #0f172a;
    padding: 12px;
    border-radius: 6px;
    border: 1px dashed #334155;
    margin-top: 14px;
  }

  .spinner {
    width: 20px;
    height: 20px;
    border: 2px solid #3b82f6;
    border-top-color: transparent;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  .spinner-msg {
    margin: 0;
    font-size: 12px;
    color: #e2e8f0;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    padding: 14px 20px;
    background: #0f172a;
  }

  .border-top {
    border-top: 1px solid #1f2937;
  }

  .input-code {
    background: #0f172a;
    border: 2px solid #3b82f6;
    border-radius: 6px;
    padding: 10px;
    width: 100%;
    color: #fff;
    font-size: 20px;
    font-weight: 700;
    letter-spacing: 4px;
    text-align: center;
    box-sizing: border-box;
  }

  /* job card details and image preview styles */
  .job-card-top {
    display: flex;
    gap: 12px;
    align-items: center;
    width: 100%;
  }

  .job-preview-img {
    width: 50px;
    height: 50px;
    object-fit: cover;
    border-radius: 4px;
    border: 1px solid #475569;
    background: #0f172a;
    flex-shrink: 0;
  }

  .job-details {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 2px;
    min-width: 0; /* allows text truncation */
  }

  .job-title-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 10px;
    width: 100%;
  }

  .job-title-text {
    font-size: 13px;
    font-weight: 700;
    color: #f8fafc;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 280px;
  }

  .btn-retry-text {
    background: transparent;
    color: #10b981;
    font-size: 12px;
    padding: 4px 8px;
    font-weight: 700;
  }

  .btn-retry-text:hover {
    background: rgba(16, 185, 129, 0.15);
    border-radius: 4px;
  }

  .btn-delete-text {
    background: transparent;
    color: #94a3b8;
    font-size: 12px;
    padding: 4px 8px;
  }

  .btn-delete-text:hover {
    background: rgba(248, 113, 113, 0.1);
    color: #fca5a5;
    border-radius: 4px;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style>
