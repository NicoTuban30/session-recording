let events = [];
let eventsURL = new URL('/events', document.currentScript.src).href;
let manualSaveTriggered = false;

rrweb.record({
  emit(event) {
    events.push(event);
  },
});

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
    agoraStreamUrl // Include the Agora stream URL
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
    events = []; // Clear events only after successful POST
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

// Optionally, save events periodically
setInterval(save, 6000); // Save every 6 seconds
