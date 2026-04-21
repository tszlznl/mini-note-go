package main

var indexHTML = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Minimalist Web Notepad</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;510;590&display=swap" rel="stylesheet">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        html, body {
            height: 100%;
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
            font-feature-settings: 'cv01', 'ss03';
        }
        body {
            display: flex;
            flex-direction: column;
            background-color: #08090a;
        }
        .header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 12px 20px;
            background-color: #0f1011;
            border-bottom: 1px solid rgba(255, 255, 255, 0.05);
            position: sticky;
            top: 0;
            z-index: 100;
        }
        .title {
            color: #f7f8f8;
            font-size: 15px;
            font-weight: 510;
            letter-spacing: -0.165px;
        }
        .actions {
            display: flex;
            gap: 8px;
        }
        .btn {
            padding: 8px 14px;
            border: none;
            border-radius: 6px;
            font-size: 13px;
            font-weight: 510;
            cursor: pointer;
            transition: all 0.2s ease;
            font-feature-settings: 'cv01', 'ss03';
        }
        .btn-primary {
            background-color: #5e6ad2;
            color: #ffffff;
        }
        .btn-primary:hover {
            background-color: #828fff;
            box-shadow: rgba(0, 0, 0, 0.1) 0px 4px 12px;
        }
        .btn-secondary {
            background-color: rgba(255, 255, 255, 0.02);
            color: #e2e4e7;
            border: 1px solid rgba(255, 255, 255, 0.08);
        }
        .btn-secondary:hover {
            background-color: rgba(255, 255, 255, 0.05);
            color: #f7f8f8;
        }
        #editor {
            flex: 1;
            width: 100%;
            padding: 20px;
            border: none;
            outline: none;
            resize: none;
            background-color: transparent;
            color: #d0d6e0;
            font-size: 16px;
            line-height: 1.5;
            font-family: inherit;
            font-feature-settings: 'cv01', 'ss03';
        }
        #editor::placeholder {
            color: #62666d;
        }
        #editor:focus {
            color: #f7f8f8;
        }
        .status-bar {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 8px 20px;
            background-color: #0f1011;
            border-top: 1px solid rgba(255, 255, 255, 0.05);
            font-size: 12px;
            color: #8a8f98;
        }
        .save-status {
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .save-indicator {
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background-color: #27a644;
            animation: pulse 2s infinite;
        }
        .save-indicator.saving {
            background-color: #ff9800;
        }
        .save-indicator.error {
            background-color: #f44336;
        }
        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.5; }
        }
        .url-display {
            font-family: ui-monospace, SF Mono, Menlo, monospace;
            font-size: 11px;
            padding: 4px 8px;
            background-color: rgba(255, 255, 255, 0.02);
            border-radius: 4px;
            border: 1px solid rgba(255, 255, 255, 0.05);
            color: #8a8f98;
            max-width: 200px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }
        @media (max-width: 600px) {
            .header {
                padding: 10px 14px;
            }
            .title {
                font-size: 14px;
            }
            .btn {
                padding: 6px 10px;
                font-size: 12px;
            }
            #editor {
                padding: 14px;
                font-size: 15px;
            }
            .status-bar {
                padding: 6px 14px;
                flex-direction: column;
                gap: 4px;
                align-items: flex-start;
            }
            .url-display {
                max-width: 100%;
            }
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="title">极简记事本</div>
        <div class="actions">
            <button class="btn btn-secondary" onclick="newNote()">新建</button>
            <button class="btn btn-primary" onclick="copyLink()">复制链接</button>
        </div>
    </div>
    <textarea id="editor" placeholder="开始输入... 您的更改将自动保存。"></textarea>
    <div class="status-bar">
        <div class="save-status">
            <span class="save-indicator" id="saveIndicator"></span>
            <span id="saveText">Saved</span>
        </div>
        <div class="url-display" id="urlDisplay"></div>
    </div>
    <script>
        const editor = document.getElementById('editor');
        const saveIndicator = document.getElementById('saveIndicator');
        const saveText = document.getElementById('saveText');
        const urlDisplay = document.getElementById('urlDisplay');
        let saveTimeout = null;
        let isSaving = false;
        let currentId = null;
        const SAVE_DELAY = 500;
        function generateId() {
            const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-';
            let id = '';
            for (let i = 0; i < 12; i++) {
                id += chars.charAt(Math.floor(Math.random() * chars.length));
            }
            return id;
        }
        function updateURL(id) {
            const newUrl = window.location.origin + '/' + id;
            window.history.pushState({ id: id }, '', newUrl);
            urlDisplay.textContent = newUrl;
            currentId = id;
        }
        function showSaveStatus(status) {
            saveIndicator.className = 'save-indicator ' + status;
            const texts = { saved: '已保存', saving: '保存中...', error: '保存失败' };
            saveText.textContent = texts[status] || '已保存';
        }
        async function saveNote() {
            if (!currentId || !editor.value.trim()) return;
            isSaving = true;
            showSaveStatus('saving');
            try {
                const response = await fetch('/' + currentId, {
                    method: 'POST',
                    body: editor.value,
                    headers: { 'Content-Type': 'text/plain' }
                });
                if (response.ok) {
                    showSaveStatus('saved');
                } else {
                    showSaveStatus('error');
                    setTimeout(() => showSaveStatus('saved'), 3000);
                }
            } catch (error) {
                showSaveStatus('error');
                setTimeout(() => showSaveStatus('saved'), 3000);
            } finally {
                isSaving = false;
            }
        }
        function scheduleSave() {
            if (saveTimeout) clearTimeout(saveTimeout);
            saveTimeout = setTimeout(saveNote, SAVE_DELAY);
        }
        async function loadNote(id) {
            try {
                const response = await fetch('/' + id + '?raw=1');
                if (response.ok) {
                    const content = await response.text();
                    editor.value = content;
                }
            } catch (error) {
                console.error('Failed to load note:', error);
            }
        }
        function newNote() {
            const id = generateId();
            updateURL(id);
            editor.value = '';
            editor.focus();
        }
        async function copyLink() {
            try {
                await navigator.clipboard.writeText(window.location.href);
                const originalText = saveText.textContent;
                saveText.textContent = '链接已复制!';
                setTimeout(() => saveText.textContent = originalText, 2000);
            } catch (error) {
                console.error('Failed to copy:', error);
            }
        }
        editor.addEventListener('input', scheduleSave);
        window.addEventListener('popstate', function(e) {
            if (e.state && e.state.id) {
                currentId = e.state.id;
                loadNote(currentId);
            }
        });
        document.addEventListener('DOMContentLoaded', function() {
            const path = window.location.pathname;
            const match = path.match(/^\/([a-zA-Z0-9_-]+)$/);
            if (match && match[1] !== 'list') {
                currentId = match[1];
                urlDisplay.textContent = window.location.href;
                loadNote(currentId);
            } else {
                newNote();
            }
            editor.focus();
        });
        window.addEventListener('beforeunload', function(e) {
            if (isSaving) {
                e.preventDefault();
                e.returnValue = '笔记保存中...';
            }
        });
    </script>
</body>
</html>
`
