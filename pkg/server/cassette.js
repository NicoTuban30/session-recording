let events = [];
let eventsURL = new URL('/events', document.currentScript.src).href;
let userEmail = getCurrentUserEmail();
let qaId = getCurrentQAUniqueId();



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

function save() {
  const body = JSON.stringify({ events, userEmail, qaId });
  events = [];

  fetch(eventsURL, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body,
  });
}

window.addEventListener('beforeunload', save);