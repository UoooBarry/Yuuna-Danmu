<script lang="ts">
  import { EventsEmit, EventsOn, Quit } from "../wailsjs/runtime/runtime";
  import { onMount, tick } from "svelte";
  import { SaveConfig, LoadConfig } from "../wailsjs/go/ui/WailsUI";
  import ToastNotification from "./lib/components/ToastNotification.svelte";
  import SuperChatCard from "./lib/components/SuperChatCard.svelte";

  let danmuList = [];
  let scList = [];
  let container: HTMLElement;
  let showControls = false;
  let showSettings = false;
  let roomID = 50819;
  let cookie = "";
  let servers: any[] = [];

  const scrollToBottom = async () => {
    if (container) {
      await tick();
      container.scrollTo({ top: container.scrollHeight, behavior: "smooth" });
    }
  };

  onMount(() => {
    LoadConfig().then(
      (config: { room_id: number; cookie: string; servers: any[] }) => {
        roomID = Number(config.room_id);
        cookie = config.cookie;
        servers = config.servers || [];
        EventsEmit("SYS_MSG", "loaded config");
      },
    );

    EventsOn("DANMU_MSG", (data) => {
      danmuList = [...danmuList, { ...data, type: "danmu" }];
      if (danmuList.length > 100) danmuList = danmuList.slice(1);
      scrollToBottom();
    });

    EventsOn("SEND_GIFT", (data) => {
      const currentComboId = data.combo_send.combo_id;

      const existingIndex = danmuList.findIndex(
        (item) => item.type === "gift" && item.combo_id === currentComboId,
      );

      let updatedGift: any;

      if (currentComboId && existingIndex !== -1) {
        const oldGift = danmuList[existingIndex];

        updatedGift = {
          ...data,
          type: "gift",
          combo_id: currentComboId,
          gift_num:
            data.combo_send?.combo_num || data.gift_num || oldGift.gift_num,
        };

        danmuList.splice(existingIndex, 1);
      } else {
        updatedGift = {
          ...data,
          type: "gift",
          combo_id: currentComboId,
          gift_num: data.combo_send?.combo_num || data.gift_num || 1,
        };
      }

      danmuList = [...danmuList, updatedGift];
      if (danmuList.length > 100) danmuList = danmuList.slice(1);
      scrollToBottom();
    });

    EventsOn("SUPER_CHAT_MESSAGE", (data) => {
      scList = [...scList, data];
    });
  });

  function handleScFinished(startTime: number) {
    scList = scList.filter((s) => s.start_time !== startTime);
  }

  async function handleSave() {
    await SaveConfig({ room_id: Number(roomID), cookie: cookie, servers: servers });
    danmuList = [];
    showSettings = false;
  }

  function getProxyUrl(originalUrl: string) {
    if (!originalUrl) return "";
    return `/proxy?url=${encodeURIComponent(originalUrl)}`;
  }

  function toggleServer(index: number) {
    servers[index].enabled = !servers[index].enabled;
    servers = [...servers]; // Trigger Svelte reactivity
  }
</script>

<main
  class="app-container"
  on:mouseenter={() => (showControls = true)}
  on:mouseleave={() => (showControls = false)}
>
  <header class="drag-bar">
    <span class="title">Yuuna Danmu</span>

    <div class="controls" class:visible={showControls}>
      <div class="controls" class:visible={showControls}>
        <button
          class="setting-btn"
          on:click={() => (showSettings = !showSettings)}
          title="设置"
        >
          <svg viewBox="0 0 24 24" width="16" height="16">
            <path
              fill="currentColor"
              d="M12,15.5A3.5,3.5 0 0,1 8.5,12A3.5,3.5 0 0,1 12,8.5A3.5,3.5 0 0,1 15.5,12A3.5,3.5 0 0,1 12,15.5M19.43,12.97C19.47,12.65 19.5,12.33 19.5,12C19.5,11.67 19.47,11.35 19.43,11.03L21.54,9.37C21.73,9.22 21.78,8.95 21.66,8.73L19.66,5.27C19.54,5.05 19.27,4.97 19.05,5.05L16.56,6.05C16.04,5.65 15.44,5.32 14.8,5.07L14.42,2.42C14.38,2.19 14.17,2 13.93,2H9.93C9.68,2 9.48,2.19 9.44,2.42L9.06,5.07C8.42,5.32 7.82,5.65 7.3,6.05L4.81,5.05C4.59,4.97 4.32,5.05 4.2,5.27L2.2,8.73C2.08,8.95 2.13,9.22 2.32,9.37L4.43,11.03C4.39,11.35 4.36,11.67 4.36,12C4.36,12.33 4.39,12.65 4.43,12.97L2.32,14.63C2.13,14.78 2.08,15.05 2.2,15.27L4.2,18.73C4.32,18.95 4.59,19.03 4.81,18.95L7.3,17.95C7.82,18.35 8.42,18.68 9.06,18.93L9.44,21.58C9.48,21.81 9.68,22 9.93,22H13.93C14.17,22 14.38,21.81 14.42,21.58L14.8,18.93C15.44,18.68 16.04,18.35 16.56,17.95L19.05,18.95C19.27,19.03 19.54,18.95 19.66,18.73L21.66,15.27C21.78,15.05 21.73,14.78 21.54,14.63L19.43,12.97Z"
            />
          </svg>
        </button>
        <button class="close-btn" on:click={Quit} title="关闭">
          <svg viewBox="0 0 24 24" width="16" height="16">
            <path
              fill="currentColor"
              d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z"
            />
          </svg>
        </button>
      </div>
    </div>
  </header>

  <div class="content-area">
    {#if showSettings}
      <div class="settings-panel">
        <div class="field">
          <label for="room-id">Room ID</label>
          <input
            type="number"
            id="room-id"
            bind:value={roomID}
            placeholder="输入直播间号"
          />
        </div>
        <div class="field">
          <label for="cookie">Cookie (SESSDATA)</label>
          <textarea
            bind:value={cookie}
            id="cookie"
            placeholder="Cookie"
            rows="5"
          ></textarea>
        </div>

        <div class="servers-section">
          <label for="servers-section">Server Settings</label>
          {#each servers as server, i}
            <div class="server-item">
              <div class="server-row">
                <span class="server-name">{server.name} ({server.type})</span>
                <label class="switch">
                  <input
                    type="checkbox"
                    checked={server.enabled}
                    on:change={() => toggleServer(i)}
                  />
                  <span class="slider"></span>
                </label>
              </div>
              <div class="server-row">
                <input
                  type="number"
                  bind:value={server.port}
                  placeholder="Port"
                />
              </div>
            </div>
          {/each}
        </div>
        <button class="save-btn" on:click={handleSave}>Apply & Restart</button>
      </div>
    {:else}
      {#if scList.length > 0}
        <div class="pinned-sc-area">
          {#each scList as sc (sc.start_time)}
            <SuperChatCard
              data={sc}
              {getProxyUrl}
              onFinished={handleScFinished}
            />
          {/each}
        </div>
      {/if}
      <div class="danmu-box" bind:this={container}>
        {#each danmuList as d}
          {#if d.type === "danmu"}
            <div class="danmu-item">
              {#if d.medalName}
                <span class="medal-tag">
                  <span class="m-name">{d.medalName}</span>
                  <span class="m-level">{d.medalLevel}</span>
                </span>
              {/if}
              <span class="nickname">{d.nickname}:</span>
              <span class="content">{d.content}</span>
            </div>
          {:else if d.type === "gift"}
            <div class="danmu-item gift-wrapper">
              <div class="gift-header">
                {#if d.face}
                  <img
                    src={getProxyUrl(d.face)}
                    alt=""
                    class="user-face-mini"
                  />
                {/if}
                {#if d.medalName}
                  <span class="medal-tag mini">
                    <span class="m-name">{d.medalName}</span>
                    <span class="m-level">{d.medalLevel}</span>
                  </span>
                {/if}
                <span class="nickname">{d.uname}</span>
                <span class="action">{d.action || "投喂"}</span>
              </div>

              <div class="gift-body">
                {#if d.gift_info}
                  <img
                    src={getProxyUrl(d.gift_info.gif)}
                    alt={d.gift_name}
                    class="gift-img"
                  />
                {:else}
                  <span class="gift-emoji">🎁</span>
                {/if}

                <div class="gift-details">
                  <span class="gift-name">{d.gift_name}</span>
                  <span class="gift-count">x {d.combo_send.combo_num}</span>
                </div>

                {#if d.coin_type === "gold"}
                  <div class="coin-badge">
                    ¥{(d.combo_total_coin / 1000).toFixed(1)}
                  </div>
                {/if}
              </div>
            </div>
          {/if}
        {/each}
      </div>
    {/if}
  </div>
</main>

<ToastNotification />

<style>
  /* Little bit of rose pine */
  .app-container {
    height: 100%;
    overflow: hidden;
    background-color: rgba(25, 23, 36, 0.7);
    backdrop-filter: blur(12px);
    border: 1px solid rgba(144, 140, 170, 0.2);
    border-radius: 12px;
    display: flex;
    flex-direction: column;
    color: var(--text);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    scrollbar-width: none;
  }

  .content-area {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
  }

  .drag-bar {
    flex-shrink: 0;
    height: 30px;
    display: flex;
    align-items: center;
    padding: 0 12px;
    font-size: 10px;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--muted);
    --wails-draggable: drag; /* 允许拖动整个窗口 */
    border-bottom: 1px solid rgba(144, 140, 170, 0.1);
  }

  .danmu-box {
    flex: 1;
    overflow-y: auto;
    padding: 10px;
    mask-image: linear-gradient(
      to bottom,
      transparent,
      black 1%,
      black 99%,
      transparent
    );
  }

  .danmu-box::-webkit-scrollbar {
    width: 4px;
  }
  .danmu-box::-webkit-scrollbar-thumb {
    border-radius: 10px;
  }

  .danmu-item {
    display: flex;
    flex-direction: row;
    justify-content: flex-start;
    align-items: flex-start;
    width: 100%;
    margin-bottom: 8px;
    text-align: left;
    word-break: break-all;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateX(5px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  .medal-tag {
    display: inline-flex;
    flex-shrink: 0;
    background: var(--overlay);
    border: 1px solid rgba(110, 106, 134, 0.5); /* --muted with alpha */
    border-radius: 4px;
    margin-right: 8px;
    height: 20px;
    align-items: center;
    overflow: hidden;
  }

  .m-name {
    display: inline-block;
    min-width: 32px;
    text-align: center;
    padding: 0 4px;
    font-size: 11px;
    color: var(--pine);
    background-color: rgba(49, 116, 143, 0.1);
  }
  .m-level {
    display: inline-block;
    min-width: 18px;
    text-align: center;
    font-size: 10px;
    color: var(--gold);
    border-left: 1px solid rgba(110, 106, 134, 0.3);
    padding: 0 3px;
  }
  .nickname {
    color: var(--rose);
    font-weight: 600;
    white-space: nowrap;
    margin-right: 6px;
    flex-shrink: 0;
  }
  .content {
    color: var(--text);
    flex: 1;
    text-align: left;
    word-break: break-all;
    white-space: pre-wrap;
  }
  .drag-bar {
    position: relative;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 8px;
    --wails-draggable: drag;
    background: rgba(31, 29, 46, 0.4);
  }

  .controls {
    opacity: 0;
    transition: opacity 0.2s ease-in-out;
    --wails-draggable: none;
    display: flex;
    align-items: center;
    display: flex;
    gap: 4px;
    opacity: 0;
  }

  .controls.visible {
    opacity: 1;
  }

  .close-btn,
  .setting-btn {
    background: transparent;
    border: none;
    color: var(--muted);
    cursor: pointer;
    padding: 4px;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
  }

  .setting-btn:hover {
    background-color: var(--pine);
    color: var(--base);
  }

  .close-btn:hover {
    background-color: var(--love);
    color: var(--base);
  }

  .title {
    font-size: 10px;
    font-weight: bold;
    color: var(--muted);
    pointer-events: none;
  }

  .settings-panel {
    padding: 15px;
    display: flex;
    flex-direction: column;
    gap: 15px;
    background: var(--surface);
    height: 100%;
    animation: slideIn 0.3s ease-out;
  }

  @keyframes slideIn {
    from {
      transform: translateY(20px);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 1;
    }
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 5px;
  }

  .field label {
    font-size: 11px;
    color: var(--iris);
    text-transform: uppercase;
  }

  input,
  textarea {
    background: var(--overlay);
    border: 1px solid var(--muted);
    color: var(--text);
    padding: 8px;
    border-radius: 6px;
    outline: none;
    font-family: inherit;
  }

  input:focus,
  textarea:focus {
    border-color: var(--pine);
  }

  .save-btn {
    background: var(--pine);
    color: var(--base);
    border: none;
    padding: 10px;
    border-radius: 6px;
    font-weight: bold;
    cursor: pointer;
    margin-top: auto;
  }

  .save-btn:hover {
    background: var(--foam);
  }

  .gift-wrapper {
    display: flex;
    flex-direction: column;
    background: linear-gradient(
      135deg,
      rgba(235, 111, 146, 0.15) 0%,
      rgba(25, 23, 36, 0.4) 100%
    );
    border-left: 3px solid var(--rose);
    border-radius: 8px;
    padding: 8px;
    margin: 6px 0;
    gap: 6px;
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
    animation: giftIn 0.4s ease-out forwards;
  }

  /* 第一行样式 */
  .gift-header {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 11px;
  }

  .user-face-mini {
    width: 18px;
    height: 18px;
    border-radius: 50%;
    border: 1px solid var(--rose);
  }

  .medal-tag.mini {
    transform: scale(0.85);
    margin-right: 0;
  }

  .action {
    color: var(--muted);
  }

  .gift-body {
    display: flex;
    align-items: center;
    width: 90%; /* 必须撑满宽度，badge 才能去到最右边 */
    gap: 10px;
    padding-top: 4px;
  }

  .gift-img {
    width: 32px;
    height: 32px;
    object-fit: contain;
    filter: drop-shadow(0 0 5px rgba(246, 193, 119, 0.5));
  }

  .gift-emoji {
    font-size: 24px;
  }

  .gift-details {
    display: flex;
    flex-direction: column;
  }

  .gift-name {
    color: var(--gold);
    font-weight: bold;
    font-size: 13px;
  }

  .gift-count {
    font-size: 16px;
    font-style: italic;
    text-shadow: 1px 1px 0 var(--base);
    display: inline-block;
    color: var(--rose);
    font-weight: 900;
    font-size: 1.4em;
    text-shadow: 2px 2px 0 var(--base);
    animation: comboPop 0.4s cubic-bezier(0.17, 0.89, 0.32, 1.49);
  }

  @keyframes comboPop {
    0% {
      transform: scale(1);
    }
    50% {
      transform: scale(1.6) rotate(-5deg);
      color: var(--gold);
    }
    100% {
      transform: scale(1) rotate(0);
    }
  }

  .coin-badge {
    margin-left: auto;
    flex-shrink: 0;

    background: rgba(246, 193, 119, 0.1);
    border: 1px solid rgba(246, 193, 119, 0.4);
    color: var(--gold);

    padding: 2px 8px;
    border-radius: 4px;
    font-size: 11px;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 4px;

    box-shadow: 0 0 8px rgba(246, 193, 119, 0.1);
  }

  .pinned-sc-area {
    display: flex;
    flex-direction: row;
    gap: 12px;
    padding: 12px;
    overflow-x: auto;
    overflow-y: hidden;
    min-height: 110px;
    background: transparent;
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);

    scrollbar-width: none;
    -ms-overflow-style: none;
  }

  .pinned-sc-area::-webkit-scrollbar {
    display: none;
  }
  .servers-section {
    color: var(--iris);
    margin-top: 1.5rem;
    font-size: 11px;
    border-top: 1px solid var(--muted);
    padding-top: 1rem;
  }

  .servers-section label {
    color: var(--iris);
    font-weight: bold;
    margin-bottom: 1rem;
  }

  .server-item {
    background: var(--overlay);
    padding: 12px;
    border-radius: 6px;
    margin-bottom: 10px;
    border: 1px solid transparent;
    transition: border 0.2s;
  }

  .server-item:hover {
    border: 1px solid var(--muted);
  }

  .server-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .server-name {
    color: var(--text);
    font-family: "JetBrains Mono", monospace; /* Very Rose Pine style */
  }

  .switch {
    position: relative;
    display: inline-block;
    width: 34px;
    height: 20px;
  }

  .switch input {
    opacity: 0;
    width: 0;
    height: 0;
  }

  .slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: var(--muted);
    transition: 0.4s;
    border-radius: 20px;
  }

  input:checked + .slider {
    background-color: var(--pine);
  }

  .slider:before {
    position: absolute;
    content: "";
    height: 14px;
    width: 14px;
    left: 3px;
    bottom: 3px;
    background-color: var(--text);
    transition: 0.4s;
    border-radius: 50%;
  }

  input:checked + .slider:before {
    transform: translateX(14px);
  }
</style>
