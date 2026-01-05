<script lang="ts">
  import { EventsOn } from '../../../wailsjs/runtime/runtime';

  interface ToastMessage {
    id: number;
    message: string;
    type: 'success' | 'error' | 'info';
  }

  let toasts: ToastMessage[] = [];
  let nextId = 0;

  EventsOn("SYS_MSG", (message: string) => {
    addToast(message, 'info');
  });

  EventsOn("SYS_ERROR", (message: string) => {
    addToast(message, 'error');
  });

  function addToast(message: string, type: ToastMessage['type'] = 'info') {
    const id = nextId++;
    toasts = [...toasts.slice(-2), { id, message, type }];
    setTimeout(() => {
      removeToast(id);
    }, 3000);
  }

  function removeToast(id: number) {
    toasts = toasts.filter(toast => toast.id !== id);
  }
</script>

<div class="toast-container">
  {#each toasts as toast (toast.id)}
    <div class="toast-item {toast.type}">
      <span class="message">{toast.message}</span>
    </div>
  {/each}
</div>

<style>
  .toast-container {
    position: fixed;
    bottom: 5px;
    left: 0;
    right: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    z-index: 9999;
    pointer-events: none;
  }

  .toast-item {
    width: 100%;
    background-color: rgba(31, 29, 46, 0.95);
    backdrop-filter: blur(8px); 
    color: #e0def4;
    padding: 6px 16px;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.3);
    border: 1px solid rgba(110, 106, 134, 0.2);
    
    animation: slideUp 0.3s ease-out, fadeOut 0.3s ease-in 2.7s forwards;
    
    white-space: nowrap;
    display: flex;
    align-items: center;
  }

  .message {
    font-size: 12px;
    font-weight: 500;
    letter-spacing: 0.3px;
  }

  .info { border-bottom: 1px solid #9ccfd8; }
  .error { border-bottom: 1px solid #eb6f92; color: #eb6f92; }

  @keyframes slideUp {
    from { opacity: 0; transform: translateY(10px) scale(0.95); }
    to { opacity: 1; transform: translateY(0) scale(1); }
  }

  @keyframes fadeOut {
    from { opacity: 1; }
    to { opacity: 0; }
  }
</style>
