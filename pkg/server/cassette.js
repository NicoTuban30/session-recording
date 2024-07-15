let events = [];
let eventsURL = new URL('/events', document.currentScript.src).href;
let userEmail = getCurrentUserEmail();
let qaId = getCurrentQAUniqueId();
let qaSessionId = getCurrentQASessionId();
let manualSaveTriggered = false;
let videoStream; // Global variable to store the video stream

rrweb.record({
  emit(event) {
    events.push(event);
  },
});

function getCurrentUserEmail() {
  const userEmail = localStorage.getItem('userEmail');
  return userEmail;
}

function getCurrentQAUniqueId() {
  const qaId = localStorage.getItem('qaId');
  return qaId;
}

function getCurrentQASessionId() {
  const qaSessionId = localStorage.getItem('qaSessionId');
  return qaSessionId;
}

function save() {
  const body = JSON.stringify({ events, userEmail, qaId, qaSessionId });
  events = [];

  fetch(eventsURL, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body,
  }).catch(error => {
    console.error('Error sending events:', error);
  });

  // If save is triggered manually, set the flag to true
  manualSaveTriggered = true;
}

window.addEventListener('beforeunload', function(event) {
  // If manualSaveTriggered is true, prevent the beforeunload event
  if (manualSaveTriggered) {
    manualSaveTriggered = false; // Reset the flag for future unload events
    return;
  }

  save();
});

