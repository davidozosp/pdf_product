document.addEventListener('DOMContentLoaded', () => {
    const fileInput = document.getElementById('fileInput');
    const fileList = document.getElementById('fileList');
    const uploadButton = document.getElementById('uploadButton');

    fileInput.addEventListener('change', () => {
        fileList.innerHTML = ''; // Clear the list
        for (const file of fileInput.files) {
            const listItem = document.createElement('li');
            listItem.textContent = file.name;
            fileList.appendChild(listItem);
        }
    });

    uploadButton.addEventListener('click', () => {
        const files = fileInput.files;
        if (files.length === 0) {
            alert('Please select files to upload.');
            return;
        }

        // Implement upload logic here
        alert('Uploading files...');
    });
});
