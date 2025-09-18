document.addEventListener('DOMContentLoaded', () => {
    // DOM Elements
    const fileUploadInput = document.getElementById('file-upload');
    const uploadButton = document.getElementById('upload-btn');
    const fileListContainer = document.getElementById('file-list');
    const logContent = document.getElementById('log-content');
    const clearLogButton = document.getElementById('clear-log-btn');

    // In-memory store for uploaded files
    let uploadedFiles = [];

    // --- Logger Function ---
    function addLog(message) {
        const timestamp = new Date().toISOString();
        const logEntry = document.createElement('div');
        logEntry.className = 'log-entry';
        logEntry.textContent = `${timestamp} - ${message}`;
        logContent.prepend(logEntry); // Add new logs to the top
    }

    // --- File List UI Renderer ---
    function renderFileList() {
        // Clear the current list
        fileListContainer.innerHTML = '';

        if (uploadedFiles.length === 0) {
            fileListContainer.innerHTML = '<p>No files uploaded yet.</p>';
            return;
        }

        // Create and append each file item
        uploadedFiles.forEach((file, index) => {
            const fileItem = document.createElement('div');
            fileItem.className = 'file-item';

            const fileName = document.createElement('span');
            fileName.textContent = file.name;

            const downloadBtn = document.createElement('button');
            downloadBtn.className = 'download-btn';
            downloadBtn.textContent = 'Download';
            
            // Event listener for downloading the file
            downloadBtn.addEventListener('click', () => {
                // Create a temporary URL for the file blob
                const url = URL.createObjectURL(file);
                
                // Create a temporary anchor tag to trigger download
                const a = document.createElement('a');
                a.style.display = 'none';
                a.href = url;
                a.download = file.name; // Set the file name for download
                
                document.body.appendChild(a);
                a.click(); // Simulate click to start download
                
                // Clean up the temporary URL and anchor tag
                URL.revokeObjectURL(url);
                document.body.removeChild(a);
                
                addLog(`Downloaded file: ${file.name}`);
            });

            fileItem.appendChild(fileName);
            fileItem.appendChild(downloadBtn);
            fileListContainer.appendChild(fileItem);
        });
    }

    // --- Event Listeners ---

    // Handle the "Upload" button click
    uploadButton.addEventListener('click', () => {
        const files = fileUploadInput.files;

        if (files.length === 0) {
            addLog('No files selected for upload.');
            return;
        }

        // Add selected files to our in-memory store
        for (const file of files) {
            // Avoid adding duplicates by checking the name and size
            if (!uploadedFiles.some(f => f.name === file.name && f.size === file.size)) {
                uploadedFiles.push(file);
                addLog(`File staged for storage: ${file.name} (${(file.size / 1024).toFixed(2)} KB)`);
            } else {
                addLog(`Skipped duplicate file: ${file.name}`);
            }
        }
        
        // Update the UI
        renderFileList();

        // Clear the file input for the next selection
        fileUploadInput.value = ''; 
    });

    // Handle the "Clear Logs" button click
    clearLogButton.addEventListener('click', () => {
        logContent.innerHTML = '';
        addLog('Logs cleared by user.');
    });

    // --- Initial State ---
    addLog('System initialized. Ready for operations.');
    renderFileList();
});