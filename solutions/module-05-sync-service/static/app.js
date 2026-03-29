const ws = new WebSocket('ws://' + location.host + '/ws');
const taskList = document.getElementById('task-list');
const statusEl = document.getElementById('status');

ws.onopen = () => {
  setStatus('connected');
  ws.send(JSON.stringify({ type: 'join' }));
};

ws.onclose = () => setStatus('disconnected');
ws.onerror = () => setStatus('disconnected');

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  if (msg.type === 'state') {
    renderTasks(msg.tasks || []);
  }
};

function setStatus(s) {
  statusEl.className = 'status ' + s;
  statusEl.textContent = s.charAt(0).toUpperCase() + s.slice(1);
}

function renderTasks(tasks) {
  taskList.innerHTML = '';
  tasks.forEach(task => {
    const li = document.createElement('li');
    if (task.completed) li.classList.add('completed');

    const checkbox = document.createElement('input');
    checkbox.type = 'checkbox';
    checkbox.checked = task.completed;
    checkbox.onchange = () => setCompleted(task.id, checkbox.checked);

    const span = document.createElement('span');
    span.textContent = task.title;
    span.title = 'Double-click to edit';
    span.ondblclick = () => editTitle(task.id, span);

    const del = document.createElement('button');
    del.textContent = '×';
    del.className = 'del-btn';
    del.onclick = () => removeTask(task.id);

    li.appendChild(checkbox);
    li.appendChild(span);
    li.appendChild(del);
    taskList.appendChild(li);
  });
}

function addTask() {
  const input = document.getElementById('new-task');
  const title = input.value.trim();
  if (!title) return;
  ws.send(JSON.stringify({ type: 'op', action: 'add', title }));
  input.value = '';
}

function removeTask(id) {
  ws.send(JSON.stringify({ type: 'op', action: 'remove', id }));
}

function setCompleted(id, completed) {
  ws.send(JSON.stringify({ type: 'op', action: 'setCompleted', id, completed }));
}

function editTitle(id, span) {
  const input = document.createElement('input');
  input.type = 'text';
  input.value = span.textContent;
  input.style.cssText = 'flex:1; padding:2px 4px; font-size:0.95rem;';
  input.onblur = () => {
    const newTitle = input.value.trim();
    if (newTitle && newTitle !== span.textContent) {
      ws.send(JSON.stringify({ type: 'op', action: 'setTitle', id, title: newTitle }));
    }
  };
  input.onkeydown = (e) => {
    if (e.key === 'Enter') input.blur();
    if (e.key === 'Escape') { input.value = span.textContent; input.blur(); }
  };
  span.replaceWith(input);
  input.focus();
  input.select();
}

document.getElementById('new-task').onkeydown = (e) => {
  if (e.key === 'Enter') addTask();
};
