let events = [];
let eventsURL = new URL('/events', document.currentScript.src).href;
let manualSaveTriggered = false;
let mediaRecorder;
let recordedBlobs;

rrweb.record({
  emit(event) {
    events.push(event);
  },
});

async function init() {
  const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
  document.querySelector('#preview').srcObject = stream;
  window.stream = stream;
  startRecording();  // Start recording automatically
}

function startRecording() {
  recordedBlobs = [];
  const options = { mimeType: 'video/webm;codecs=vp9' };
  mediaRecorder = new MediaRecorder(window.stream, options);
  mediaRecorder.ondataavailable = handleDataAvailable;
  mediaRecorder.start();
}

function stopRecording() {
  mediaRecorder.stop();
  const superBuffer = new Blob(recordedBlobs, { type: 'video/webm' });
  document.querySelector('#recorded').src = window.URL.createObjectURL(superBuffer);

  const formData = new FormData();
  formData.append('session', getCurrentQASessionId());
  formData.append('video', superBuffer, 'video.webm');

  fetch('/upload', {
    method: 'POST',
    body: formData
  });
}

function handleDataAvailable(event) {
  if (event.data && event.data.size > 0) {
    recordedBlobs.push(event.data);
  }
}

document.querySelector('#start').addEventListener('click', startRecording);
document.querySelector('#stop').addEventListener('click', stopRecording);

init();

function getCurrentUserEmail() {
  return localStorage.getItem('userEmail') || '';
}

function getCurrentQAUniqueId() {
  return localStorage.getItem('qaId') || '';
}

function getCurrentQASessionId() {
  return localStorage.getItem('qaSessionId') || '';
}

function getCurrentAgoraStreamUrl() {
  return localStorage.getItem('agoraStreamUrl') || '';
}

function save() {
  let userEmail = getCurrentUserEmail();
  let qaId = getCurrentQAUniqueId();
  let qaSessionId = getCurrentQASessionId();
  let agoraStreamUrl = getCurrentAgoraStreamUrl();

  if (!qaId || !qaSessionId) {
    console.warn('Missing userEmail, qaId, or qaSessionId. Events not saved.');
    return;
  }

  const body = JSON.stringify({
    events,
    userEmail,
    qaId,
    qaSessionId,
    agoraStreamUrl
  });

  fetch(eventsURL, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body,
  }).then(response => {
    if (!response.ok) {
      return Promise.reject('Failed to save events');
    }
    events = [];
  }).catch(error => {
    console.error('Error sending events:', error);
  });

  manualSaveTriggered = true;
}

window.addEventListener('beforeunload', function(event) {
  if (manualSaveTriggered) {
    manualSaveTriggered = false;
    return;
  }
  save();
});

setInterval(save, 6000);
