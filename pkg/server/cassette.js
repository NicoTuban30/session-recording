let events = [];
let eventsURL = new URL('/events', document.currentScript.src).href;
let userId = getCurrentUserId();


rrweb.record({
  emit(event) {
    events.push(event);
  },
});

function getCurrentUserId() {
  const userId = localStorage.getItem('userId');
  
  return userId;

}

function save() {
  const body = JSON.stringify({ events, userId });
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