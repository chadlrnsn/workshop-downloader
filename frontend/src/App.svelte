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

  // Subcomponents import
  import Header from './components/Header.svelte';
  import NewTaskCard from './components/NewTaskCard.svelte';
  import StatsCard from './components/StatsCard.svelte';
  import SessionCard from './components/SessionCard.svelte';
  import QueuePanel from './components/QueuePanel.svelte';
  import ConsolePanel from './components/ConsolePanel.svelte';
  import SettingsModal from './components/SettingsModal.svelte';
  import AuthCodeModal from './components/AuthCodeModal.svelte';

  // Component State
  let config = { steamCmdPath: '', outputDir: '', autoUpdate: true, username: 'anonymous' };
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

  async function handleSubmitDownload(event) {
    const { workshopUrl, manualAppId } = event.detail;
    try {
      await AddDownload(workshopUrl, manualAppId);
      await fetchJobs();
    } catch (err) {
      alert(`Error submitting job: ${err}`);
    }
  }

  async function handleCancel(event) {
    const jobId = event.detail;
    await CancelJob(jobId);
    await fetchJobs();
  }

  async function handleRetry(event) {
    const jobId = event.detail;
    try {
      await RetryJob(jobId);
      await fetchJobs();
    } catch (err) {
      alert(`Retry failed: ${err}`);
    }
  }

  async function handleDelete(event) {
    const jobId = event.detail;
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
  <!-- Top Bar Custom Component -->
  <Header on:toggleSettings={() => showConfigModal = true} />

  <div class="workspace-body">
    <!-- Left Column: Input Forms & Stats Component Pane -->
    <aside class="left-pane">
      <NewTaskCard on:submit={handleSubmitDownload} />
      <StatsCard {jobs} />
      <SessionCard username={config.username} />
    </aside>

    <!-- Right Column: Queue & Terminal Panel Components -->
    <section class="right-pane">
      <QueuePanel 
        {jobs} 
        on:cancel={handleCancel}
        on:retry={handleRetry}
        on:delete={handleDelete}
      />
      <ConsolePanel 
        {activeJobLogs} 
        bind:logContainer 
        on:clear={() => activeJobLogs = []} 
      />
    </section>
  </div>

  <!-- Settings & Authentication Modal Overlay -->
  {#if showConfigModal}
    <SettingsModal 
      bind:config 
      bind:steamPassword 
      {isSavingSettings}
      {isLoggingIn}
      {isCheckingSteamCmd}
      {loginStatusMsg}
      on:save={updateConfig}
      on:check={handleVerifySteamCmd}
      on:login={triggerSteamLogin}
      on:close={() => { showConfigModal = false; steamPassword = ''; }}
    />
  {/if}

  <!-- Modal Dialog: Steam Guard Authenticator -->
  {#if showAuthModal}
    <AuthCodeModal 
      {authPromptType}
      bind:steamCode 
      on:submit={submitCode}
      on:cancel={cancelCode}
    />
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

  /* Fullpage Layout Core */
  .app-workspace {
    display: flex;
    flex-direction: column;
    width: 100vw;
    height: 100vh;
    box-sizing: border-box;
  }

  .workspace-body {
    display: flex;
    flex: 1;
    height: calc(100vh - 40px);
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

  /* Global buttons shared parameters fallback */
  :global(button) {
    font-family: inherit;
    font-size: 13px;
    font-weight: 600;
    border-radius: 6px;
    padding: 8px 16px;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    border: none;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
  }

  :global(button:disabled) {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  :global(.btn-primary) {
    background: #2563eb;
    color: #ffffff;
  }
  :global(.btn-primary:hover:not(:disabled)) {
    background: #1d4ed8;
  }
  :global(.btn-primary.btn-block) {
    width: 100%;
    display: flex;
    justify-content: center;
    padding: 11px;
    box-sizing: border-box;
  }
  :global(.btn-secondary) {
    background: #1e293b;
    border: 1px solid #334155;
    color: #cbd5e1;
  }
  :global(.btn-secondary:hover:not(:disabled)) {
    background: #334155;
    color: #f8fafc;
  }
</style>
