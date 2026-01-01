package auth

const setupTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Control D CLI - Connect Account</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg: #010818;
            --bg-card: #0a1628;
            --bg-input: #131f35;
            --bg-hint: rgba(74, 32, 229, 0.08);
            --border: #1e3a5f;
            --border-focus: #4a20e5;
            --text: #f3f4f6;
            --text-secondary: #a78fff;
            --text-muted: #6b7280;
            --primary: #4a20e5;
            --primary-light: rgba(74, 32, 229, 0.15);
            --accent: #a78fff;
            --success: #10b981;
            --success-light: rgba(16, 185, 129, 0.15);
            --error: #ef4444;
            --error-light: rgba(239, 68, 68, 0.15);
            --glow: 0 0 20px rgba(74, 32, 229, 0.4);
        }

        * { margin: 0; padding: 0; box-sizing: border-box; }

        html { height: 100%; }

        body {
            font-family: 'Plus Jakarta Sans', -apple-system, sans-serif;
            background: var(--bg);
            color: var(--text);
            min-height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 2rem 1.5rem;
            background-image:
                radial-gradient(ellipse at 50% 0%, rgba(74, 32, 229, 0.12) 0%, transparent 50%),
                radial-gradient(circle at 80% 80%, rgba(167, 143, 255, 0.06) 0%, transparent 40%);
        }

        .container {
            width: 100%;
            max-width: 400px;
        }

        .logo {
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 0.75rem;
            margin-bottom: 1.5rem;
        }

        .logo-icon {
            width: 40px;
            height: 40px;
            background: linear-gradient(135deg, var(--primary) 0%, #7c3aed 100%);
            border-radius: 10px;
            display: flex;
            align-items: center;
            justify-content: center;
            box-shadow: var(--glow);
        }

        .logo-icon svg {
            width: 24px;
            height: 24px;
            color: white;
        }

        .logo-text {
            font-size: 1.5rem;
            font-weight: 700;
            letter-spacing: -0.02em;
            background: linear-gradient(135deg, var(--text) 0%, var(--text-secondary) 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }

        .badge-wrapper {
            display: flex;
            justify-content: center;
            margin-bottom: 1.25rem;
        }

        .cli-badge {
            display: inline-flex;
            align-items: center;
            gap: 0.375rem;
            background: var(--primary-light);
            color: var(--primary);
            font-size: 0.6875rem;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.08em;
            padding: 0.375rem 0.75rem;
            border-radius: 100px;
            border: 1px solid rgba(6, 182, 212, 0.2);
        }

        .cli-badge svg {
            width: 12px;
            height: 12px;
        }

        h1 {
            font-size: 1.375rem;
            font-weight: 600;
            text-align: center;
            margin-bottom: 0.375rem;
        }

        .subtitle {
            color: var(--text-secondary);
            font-size: 0.875rem;
            text-align: center;
            margin-bottom: 1.25rem;
        }

        .credentials-hint {
            background: var(--bg-hint);
            border: 1px solid rgba(6, 182, 212, 0.15);
            border-radius: 10px;
            padding: 0.875rem;
            margin-bottom: 1rem;
        }

        .hint-header {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            font-size: 0.75rem;
            font-weight: 500;
            color: var(--primary);
            margin-bottom: 0.625rem;
        }

        .hint-header svg {
            width: 14px;
            height: 14px;
        }

        .hint-link {
            display: flex;
            align-items: center;
            gap: 0.625rem;
            padding: 0.625rem 0.75rem;
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 8px;
            text-decoration: none;
            color: var(--text);
            transition: all 0.15s ease;
        }

        .hint-link:hover {
            border-color: var(--primary);
            box-shadow: 0 0 0 2px rgba(6, 182, 212, 0.1);
        }

        .hint-link-icon {
            width: 28px;
            height: 28px;
            background: var(--primary-light);
            border-radius: 6px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: var(--primary);
            flex-shrink: 0;
        }

        .hint-link-icon svg {
            width: 14px;
            height: 14px;
        }

        .hint-link-text {
            flex: 1;
            min-width: 0;
        }

        .hint-link-title {
            font-weight: 600;
            font-size: 0.8125rem;
        }

        .hint-link-path {
            font-size: 0.6875rem;
            color: var(--text-muted);
            font-family: 'IBM Plex Mono', monospace;
        }

        .hint-link-arrow {
            color: var(--text-muted);
            flex-shrink: 0;
        }

        .form-card {
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 1.25rem;
        }

        .form-group {
            margin-bottom: 1rem;
        }

        .form-group:last-of-type {
            margin-bottom: 0;
        }

        .label-row {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 0.375rem;
        }

        label {
            font-size: 0.8125rem;
            font-weight: 600;
            color: var(--text);
        }

        .badge {
            font-size: 0.5625rem;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.06em;
            padding: 0.1875rem 0.4375rem;
            border-radius: 4px;
        }

        .badge-required {
            background: var(--primary-light);
            color: var(--primary);
        }

        .input-wrapper {
            position: relative;
        }

        input {
            width: 100%;
            padding: 0.625rem 0.875rem;
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.8125rem;
            background: var(--bg-input);
            border: 1.5px solid var(--border);
            border-radius: 8px;
            color: var(--text);
            transition: all 0.15s ease;
        }

        input::placeholder {
            color: var(--text-muted);
        }

        input:focus {
            outline: none;
            border-color: var(--primary);
            box-shadow: 0 0 0 3px rgba(6, 182, 212, 0.15);
        }

        input.error {
            border-color: var(--error);
            background: var(--error-light);
        }

        input.error:focus {
            border-color: var(--error);
            box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.15);
        }

        .input-hint {
            font-size: 0.6875rem;
            color: var(--text-muted);
            margin-top: 0.25rem;
        }

        .password-toggle {
            position: absolute;
            right: 0.5rem;
            top: 50%;
            transform: translateY(-50%);
            background: none;
            border: none;
            color: var(--text-muted);
            cursor: pointer;
            padding: 0.25rem;
            border-radius: 4px;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .password-toggle:hover {
            color: var(--text-secondary);
        }

        .password-toggle svg {
            width: 16px;
            height: 16px;
        }

        .btn-group {
            display: flex;
            gap: 0.625rem;
            margin-top: 1.25rem;
        }

        button {
            flex: 1;
            padding: 0.6875rem 1rem;
            font-family: inherit;
            font-size: 0.8125rem;
            font-weight: 600;
            border-radius: 8px;
            cursor: pointer;
            transition: all 0.15s ease;
            border: none;
        }

        .btn-secondary {
            background: var(--bg-input);
            color: var(--text-secondary);
            border: 1px solid var(--border);
        }

        .btn-secondary:hover {
            background: var(--border);
            color: var(--text);
        }

        .btn-primary {
            background: linear-gradient(135deg, var(--primary) 0%, #7c3aed 100%);
            color: white;
            box-shadow: 0 2px 8px rgba(74, 32, 229, 0.3);
        }

        .btn-primary:hover {
            box-shadow: 0 4px 16px rgba(74, 32, 229, 0.4);
            transform: translateY(-1px);
        }

        button:disabled {
            opacity: 0.5;
            cursor: not-allowed;
            transform: none;
        }

        .status {
            position: fixed;
            bottom: 1.5rem;
            left: 50%;
            transform: translateX(-50%) translateY(10px);
            padding: 0.625rem 1rem;
            border-radius: 8px;
            font-size: 0.75rem;
            font-weight: 500;
            align-items: center;
            gap: 0.5rem;
            opacity: 0;
            visibility: hidden;
            transition: all 0.2s ease;
            display: flex;
            box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
            z-index: 100;
            white-space: nowrap;
            border: 1px solid transparent;
        }

        .status.show {
            opacity: 1;
            visibility: visible;
            transform: translateX(-50%) translateY(0);
        }

        .status.loading {
            background: var(--bg-card);
            color: var(--primary);
            border-color: rgba(6, 182, 212, 0.3);
        }

        .status.success {
            background: var(--success-light);
            color: var(--success);
            border-color: rgba(16, 185, 129, 0.3);
        }

        .status.error {
            background: var(--error-light);
            color: var(--error);
            border-color: rgba(239, 68, 68, 0.3);
        }

        .spinner {
            width: 14px;
            height: 14px;
            border: 2px solid currentColor;
            border-top-color: transparent;
            border-radius: 50%;
            animation: spin 0.6s linear infinite;
        }

        @keyframes spin { to { transform: rotate(360deg); } }

        .status-icon { width: 14px; height: 14px; flex-shrink: 0; }

        .github-link {
            position: fixed;
            bottom: 1.5rem;
            left: 50%;
            transform: translateX(-50%);
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
            text-decoration: none;
            color: var(--text-muted);
            font-size: 0.75rem;
            font-weight: 500;
            transition: color 0.2s ease;
        }

        .github-link:hover {
            color: var(--accent);
        }

        .github-icon {
            width: 16px;
            height: 16px;
        }

        .logo-svg {
            height: 48px;
            width: auto;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">
            <svg class="logo-svg" xmlns="http://www.w3.org/2000/svg" width="42" height="38" viewBox="0 0 42 38" fill="none">
                <path d="M2.59817 16C1.73199 16 1.03883 15.2655 1.18488 14.4118C2.39502 7.338 8.52232 2 16.2388 2H19.3786C20.0414 2 20.5786 2.53726 20.5786 3.2C20.5786 10.2692 14.8479 16 7.77861 16H2.59817Z" fill="#8B92A2"/>
                <path d="M2.43556 17.9426C1.56751 17.9398 0.871345 18.6752 1.02014 19.5304C2.25182 26.6097 8.52372 32 16.2388 32H19.3751C20.0398 32 20.5786 31.4612 20.5786 30.7965C20.5786 23.7227 14.856 17.9819 7.78215 17.9595L2.43556 17.9426Z" fill="#8B92A2"/>
                <path d="M35.3868 18C28.3175 18 22.5868 23.7308 22.5868 30.8C22.5868 31.4627 23.124 32 23.7868 32H26.9273C34.6437 32 40.7705 26.662 41.9805 19.5882C42.1266 18.7345 41.4334 18 40.5672 18H35.3868Z" fill="#8B92A2"/>
                <path d="M40.5672 16C41.4334 16 42.1266 15.2655 41.9805 14.4118C40.7705 7.338 34.6437 2 26.9273 2H23.7135C23.0444 2 22.5031 2.54462 22.5073 3.21371C22.5513 10.2882 28.2987 16 35.3733 16H40.5672Z" fill="#8B92A2"/>
            </svg>
        </div>

        <div class="badge-wrapper">
            <div class="cli-badge">
                <svg viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <rect x="2" y="2" width="12" height="12" rx="2" stroke="currentColor" stroke-width="1.5"/>
                    <path d="M5 6L7 8L5 10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M9 10H11" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                </svg>
                CLI Authentication
            </div>
        </div>

        <h1>Connect Your Account</h1>
        <p class="subtitle">Enter your Control D API token to get started</p>

        <div class="credentials-hint">
            <div class="hint-header">
                <svg viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M8 2v4M8 10v4M2 8h4M10 8h4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                </svg>
                Get your API token
            </div>
            <a href="https://controld.com/dashboard/settings" target="_blank" class="hint-link">
                <div class="hint-link-icon">
                    <svg viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <rect x="3" y="7" width="10" height="6" rx="1" stroke="currentColor" stroke-width="1.5"/>
                        <path d="M5 7V5a3 3 0 0 1 6 0v2" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                        <circle cx="8" cy="10" r="1" fill="currentColor"/>
                    </svg>
                </div>
                <div class="hint-link-text">
                    <div class="hint-link-title">Control D Dashboard</div>
                    <div class="hint-link-path">Settings &rarr; API</div>
                </div>
                <svg class="hint-link-arrow" width="16" height="16" viewBox="0 0 16 16" fill="none">
                    <path d="M6 4L10 8L6 12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
            </a>
        </div>

        <div class="form-card">
            <form id="setupForm" autocomplete="off">
                <div class="form-group">
                    <div class="label-row">
                        <label for="accountName">Account Name</label>
                        <span class="badge badge-required">Required</span>
                    </div>
                    <input type="text" id="accountName" name="accountName" placeholder="e.g., personal, work" required autofocus>
                    <div class="input-hint">Name can be anything - it's just a local identifier</div>
                </div>

                <div class="form-group">
                    <div class="label-row">
                        <label for="apiToken">API Token</label>
                        <span class="badge badge-required">Required</span>
                    </div>
                    <div class="input-wrapper">
                        <input type="password" id="apiToken" name="apiToken" placeholder="api.xxxxxxxxxxxxxxxxxxxxxxxx" required style="padding-right: 2rem;">
                        <button type="button" class="password-toggle" id="togglePassword" aria-label="Toggle visibility">
                            <svg id="eyeIcon" viewBox="0 0 18 18" fill="none">
                                <path d="M2 9C2 9 4.5 4 9 4C13.5 4 16 9 16 9C16 9 13.5 14 9 14C4.5 14 2 9 2 9Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                                <circle cx="9" cy="9" r="2" stroke="currentColor" stroke-width="1.5"/>
                            </svg>
                            <svg id="eyeOffIcon" style="display:none" viewBox="0 0 18 18" fill="none">
                                <path d="M7.6 7.6a2 2 0 1 0 2.8 2.8M12.5 12.5A6.5 6.5 0 0 1 9 14c-4.5 0-7-5-7-5a11.5 11.5 0 0 1 3-3.5m2.2-1.2A5.5 5.5 0 0 1 9 4c4.5 0 7 5 7 5a11.5 11.5 0 0 1-1.2 1.8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                                <path d="M2 2l14 14" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                            </svg>
                        </button>
                    </div>
                    <div class="input-hint">Starts with api. followed by your token</div>
                </div>

                <div class="btn-group">
                    <button type="button" id="testBtn" class="btn-secondary">Test</button>
                    <button type="submit" id="submitBtn" class="btn-primary">Save & Connect</button>
                </div>

                <div id="status" class="status"></div>
            </form>
        </div>
    </div>

    <a href="https://github.com/salmonumbrella/controld-cli" target="_blank" class="github-link">
        <svg class="github-icon" viewBox="0 0 16 16" fill="currentColor">
            <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
        </svg>
        View on GitHub
    </a>

    <script>
        const form = document.getElementById('setupForm');
        const testBtn = document.getElementById('testBtn');
        const submitBtn = document.getElementById('submitBtn');
        const status = document.getElementById('status');
        const togglePassword = document.getElementById('togglePassword');
        const apiTokenInput = document.getElementById('apiToken');
        const eyeIcon = document.getElementById('eyeIcon');
        const eyeOffIcon = document.getElementById('eyeOffIcon');
        const csrfToken = '{{.CSRFToken}}';

        const requiredFields = ['accountName', 'apiToken'];
        let isBusy = false;

        requiredFields.forEach(id => {
            document.getElementById(id).addEventListener('input', function() {
                this.classList.remove('error');
            });
        });

        togglePassword.addEventListener('click', () => {
            const isPassword = apiTokenInput.type === 'password';
            apiTokenInput.type = isPassword ? 'text' : 'password';
            eyeIcon.style.display = isPassword ? 'none' : 'block';
            eyeOffIcon.style.display = isPassword ? 'block' : 'none';
        });

        function showStatus(type, message) {
            status.className = 'status show ' + type;
            if (type === 'loading') {
                status.innerHTML = '<div class="spinner"></div><span>' + message + '</span>';
            } else {
                const icon = type === 'success'
                    ? '<svg class="status-icon" viewBox="0 0 16 16" fill="none"><path d="M13 5L6.5 11.5L3 8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>'
                    : '<svg class="status-icon" viewBox="0 0 16 16" fill="none"><path d="M12 4L4 12M4 4L12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>';
                status.innerHTML = icon + '<span>' + message + '</span>';
            }
        }

        function hideStatus() {
            status.className = 'status';
        }

        function validateRequired() {
            let valid = true;
            requiredFields.forEach(id => {
                const input = document.getElementById(id);
                if (!input.value.trim()) {
                    input.classList.add('error');
                    valid = false;
                }
            });
            return valid;
        }

        function getFormData() {
            return {
                account_name: document.getElementById('accountName').value.trim(),
                api_token: document.getElementById('apiToken').value.trim()
            };
        }

        testBtn.addEventListener('click', async () => {
            if (isBusy) return;
            isBusy = true;
            hideStatus();
            if (!validateRequired()) {
                isBusy = false;
                return;
            }

            const data = getFormData();
            testBtn.disabled = true;
            submitBtn.disabled = true;
            showStatus('loading', 'Testing connection...');
            try {
                const response = await fetch('/validate', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json', 'X-CSRF-Token': csrfToken },
                    body: JSON.stringify(data)
                });
                const result = await response.json();
                showStatus(result.success ? 'success' : 'error', result.success ? 'Connection successful!' : result.error);
            } catch (err) {
                showStatus('error', 'Request failed: ' + err.message);
            } finally {
                testBtn.disabled = false;
                submitBtn.disabled = false;
                isBusy = false;
            }
        });

        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            if (isBusy) return;
            isBusy = true;
            hideStatus();
            if (!validateRequired()) {
                isBusy = false;
                return;
            }

            const data = getFormData();
            testBtn.disabled = true;
            submitBtn.disabled = true;
            showStatus('loading', 'Saving credentials...');
            try {
                const response = await fetch('/submit', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json', 'X-CSRF-Token': csrfToken },
                    body: JSON.stringify(data)
                });
                const result = await response.json();
                if (result.success) {
                    showStatus('success', 'Credentials saved! Redirecting...');
                    setTimeout(() => { window.location.href = '/success'; }, 600);
                } else {
                    showStatus('error', result.error);
                    testBtn.disabled = false;
                    submitBtn.disabled = false;
                    isBusy = false;
                }
            } catch (err) {
                showStatus('error', 'Request failed: ' + err.message);
                testBtn.disabled = false;
                submitBtn.disabled = false;
                isBusy = false;
            }
        });
    </script>
</body>
</html>`

const successTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Connected - Control D CLI</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg: #010818;
            --bg-card: #0a1628;
            --bg-terminal: #050d1a;
            --border: #1e3a5f;
            --text: #f3f4f6;
            --text-secondary: #a78fff;
            --text-muted: #6b7280;
            --primary: #4a20e5;
            --primary-light: rgba(74, 32, 229, 0.15);
            --accent: #a78fff;
            --success: #10b981;
            --success-light: rgba(16, 185, 129, 0.15);
        }

        * { margin: 0; padding: 0; box-sizing: border-box; }
        html { height: 100%; }

        body {
            font-family: 'Plus Jakarta Sans', -apple-system, sans-serif;
            background: var(--bg);
            color: var(--text);
            min-height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 2rem 1.5rem;
            background-image:
                radial-gradient(ellipse at 50% 0%, rgba(16, 185, 129, 0.1) 0%, transparent 50%),
                radial-gradient(circle at 20% 80%, rgba(74, 32, 229, 0.08) 0%, transparent 40%);
        }

        .container {
            width: 100%;
            max-width: 400px;
            text-align: center;
        }

        .success-icon {
            width: 64px;
            height: 64px;
            background: var(--success-light);
            border: 2px solid rgba(16, 185, 129, 0.3);
            border-radius: 50%;
            margin: 0 auto 1.25rem;
            display: flex;
            align-items: center;
            justify-content: center;
            animation: scaleIn 0.5s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
            box-shadow: 0 0 30px rgba(16, 185, 129, 0.2);
        }

        @keyframes scaleIn {
            from { transform: scale(0); opacity: 0; }
            to { transform: scale(1); opacity: 1; }
        }

        .success-icon svg {
            width: 32px;
            height: 32px;
            color: var(--success);
        }

        h1 {
            font-size: 1.5rem;
            font-weight: 700;
            margin-bottom: 0.375rem;
            animation: fadeUp 0.5s ease 0.2s both;
        }

        .subtitle {
            color: var(--text-secondary);
            font-size: 0.875rem;
            margin-bottom: 1rem;
            animation: fadeUp 0.5s ease 0.3s both;
        }

        @keyframes fadeUp {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }

        .account-badge {
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
            background: var(--primary-light);
            border: 1px solid rgba(74, 32, 229, 0.2);
            color: var(--accent);
            font-size: 0.8125rem;
            font-weight: 600;
            padding: 0.5rem 1rem;
            border-radius: 100px;
            margin-bottom: 1.5rem;
            animation: fadeUp 0.5s ease 0.35s both;
        }

        .account-badge .dot {
            width: 8px;
            height: 8px;
            background: var(--success);
            border-radius: 50%;
            animation: dotPulse 2s ease-in-out infinite;
            box-shadow: 0 0 8px rgba(16, 185, 129, 0.5);
        }

        @keyframes dotPulse {
            0%, 100% { opacity: 1; transform: scale(1); }
            50% { opacity: 0.6; transform: scale(0.9); }
        }

        .terminal {
            background: var(--bg-terminal);
            border: 1px solid var(--border);
            border-radius: 10px;
            overflow: hidden;
            text-align: left;
            animation: fadeUp 0.5s ease 0.4s both;
        }

        .terminal-bar {
            background: #161b22;
            padding: 0.625rem 0.875rem;
            display: flex;
            align-items: center;
            gap: 0.5rem;
            border-bottom: 1px solid var(--border);
        }

        .terminal-dot {
            width: 10px;
            height: 10px;
            border-radius: 50%;
        }

        .terminal-dot.red { background: #ff5f57; }
        .terminal-dot.yellow { background: #febc2e; }
        .terminal-dot.green { background: #28c840; }

        .terminal-body {
            padding: 1rem;
        }

        .terminal-line {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.75rem;
            margin-bottom: 0.625rem;
            color: #e6edf3;
        }

        .terminal-line:last-child { margin-bottom: 0; }
        .terminal-prompt { color: var(--primary); user-select: none; }
        .terminal-cmd { color: var(--success); }
        .terminal-output {
            color: var(--text-muted);
            padding-left: 1.25rem;
            margin-top: -0.375rem;
            margin-bottom: 0.625rem;
            font-size: 0.6875rem;
        }

        .terminal-cursor {
            display: inline-block;
            width: 8px;
            height: 16px;
            background: var(--primary);
            animation: cursorBlink 1.2s step-end infinite;
            margin-left: 2px;
            vertical-align: middle;
        }

        @keyframes cursorBlink {
            0%, 50% { opacity: 1; }
            50.01%, 100% { opacity: 0; }
        }

        .message {
            margin-top: 1.25rem;
            padding: 1rem;
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 10px;
            text-align: center;
            animation: fadeUp 0.5s ease 0.5s both;
        }

        .message-icon {
            font-size: 1.25rem;
            margin-bottom: 0.25rem;
            color: var(--primary);
        }

        .message-title {
            font-weight: 600;
            font-size: 0.875rem;
            margin-bottom: 0.25rem;
        }

        .message-text {
            font-size: 0.75rem;
            color: var(--text-secondary);
            line-height: 1.5;
        }

        .message-text code {
            font-family: 'JetBrains Mono', monospace;
            background: rgba(74, 32, 229, 0.15);
            color: var(--accent);
            padding: 0.125rem 0.375rem;
            border-radius: 4px;
            font-size: 0.6875rem;
        }

        .footer {
            margin-top: 1rem;
            font-size: 0.6875rem;
            color: var(--text-muted);
            animation: fadeUp 0.5s ease 0.6s both;
        }

        .github-link {
            position: fixed;
            bottom: 1.5rem;
            left: 50%;
            transform: translateX(-50%);
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
            text-decoration: none;
            color: var(--text-muted);
            font-size: 0.75rem;
            font-weight: 500;
            transition: color 0.2s ease;
            animation: fadeUp 0.5s ease 0.7s both;
        }

        .github-link:hover {
            color: var(--accent);
        }

        .github-icon {
            width: 16px;
            height: 16px;
        }

        .logo {
            display: flex;
            justify-content: center;
            margin-bottom: 1.5rem;
            animation: fadeUp 0.5s ease 0.1s both;
        }

        .logo-svg {
            height: 48px;
            width: auto;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">
            <svg class="logo-svg" xmlns="http://www.w3.org/2000/svg" width="42" height="38" viewBox="0 0 42 38" fill="none">
                <path d="M2.59817 16C1.73199 16 1.03883 15.2655 1.18488 14.4118C2.39502 7.338 8.52232 2 16.2388 2H19.3786C20.0414 2 20.5786 2.53726 20.5786 3.2C20.5786 10.2692 14.8479 16 7.77861 16H2.59817Z" fill="#8B92A2"/>
                <path d="M2.43556 17.9426C1.56751 17.9398 0.871345 18.6752 1.02014 19.5304C2.25182 26.6097 8.52372 32 16.2388 32H19.3751C20.0398 32 20.5786 31.4612 20.5786 30.7965C20.5786 23.7227 14.856 17.9819 7.78215 17.9595L2.43556 17.9426Z" fill="#8B92A2"/>
                <path d="M35.3868 18C28.3175 18 22.5868 23.7308 22.5868 30.8C22.5868 31.4627 23.124 32 23.7868 32H26.9273C34.6437 32 40.7705 26.662 41.9805 19.5882C42.1266 18.7345 41.4334 18 40.5672 18H35.3868Z" fill="#8B92A2"/>
                <path d="M40.5672 16C41.4334 16 42.1266 15.2655 41.9805 14.4118C40.7705 7.338 34.6437 2 26.9273 2H23.7135C23.0444 2 22.5031 2.54462 22.5073 3.21371C22.5513 10.2882 28.2987 16 35.3733 16H40.5672Z" fill="#8B92A2"/>
            </svg>
        </div>

        <div class="success-icon">
            <svg viewBox="0 0 32 32" fill="none">
                <path d="M8 16L14 22L24 10" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
        </div>

        <h1>You're all set!</h1>
        <p class="subtitle">Control D CLI is now connected and ready</p>

        {{if .AccountName}}
        <div class="account-badge">
            <span class="dot"></span>
            <span>{{.AccountName}}</span>
        </div>
        {{end}}

        <div class="terminal">
            <div class="terminal-bar">
                <span class="terminal-dot red"></span>
                <span class="terminal-dot yellow"></span>
                <span class="terminal-dot green"></span>
            </div>
            <div class="terminal-body">
                <div class="terminal-line">
                    <span class="terminal-prompt">$</span>
                    <span class="terminal-cmd">controld</span>
                    <span>profiles list</span>
                </div>
                <div class="terminal-output">Fetching DNS profiles...</div>
                <div class="terminal-line">
                    <span class="terminal-prompt">$</span>
                    <span class="terminal-cmd">controld</span>
                    <span>devices list</span>
                </div>
                <div class="terminal-output">Listing managed devices...</div>
                <div class="terminal-line">
                    <span class="terminal-prompt">$</span>
                    <span class="terminal-cursor"></span>
                </div>
            </div>
        </div>

        <div class="message">
            <div class="message-icon">&larr;</div>
            <div class="message-title">Return to your terminal</div>
            <div class="message-text">Close this window and start using the CLI.<br>Try <code>controld --help</code> to see all commands.</div>
        </div>

        <p class="footer">This window will close automatically.</p>
    </div>

    <a href="https://github.com/salmonumbrella/controld-cli" target="_blank" class="github-link">
        <svg class="github-icon" viewBox="0 0 16 16" fill="currentColor">
            <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
        </svg>
        View on GitHub
    </a>

    <script>fetch('/complete', { method: 'POST', headers: { 'X-CSRF-Token': '{{.CSRFToken}}' } }).catch(() => {});</script>
</body>
</html>`
