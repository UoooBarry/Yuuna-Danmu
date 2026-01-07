<script lang="ts">
  import { onMount, onDestroy } from "svelte";

  export let data: any;
  export let getProxyUrl: (url: string) => string;
  export let onFinished: (startTime: number) => void;

  let remainingPercent = 100;
  let timer: any;

  function updateTimer() {
    const now = Math.floor(Date.now() / 1000);
    const total = data.end_time - data.start_time;
    const remaining = data.end_time - now;

    if (remaining <= 0) {
      remainingPercent = 0;
      clearInterval(timer);
      onFinished(data.start_time);
    } else {
      remainingPercent = (remaining / total) * 100;
    }
  }

  onMount(() => {
    updateTimer();
    timer = setInterval(updateTimer, 1000);
  });

  onDestroy(() => {
    clearInterval(timer);
  });

  const getScColor = (price: number) => {
    if (price >= 1000) return "#eb6f92"; // love (红)
    if (price >= 500) return "#c4a7e7";  // iris (紫)
    if (price >= 100) return "#f6c177";  // gold (金)
    return "#31748f";                    // pine (青)
  };
</script>

<div class="sc-card" style="--sc-bg: {getScColor(data.price)}">
  <div class="sc-header">
    <img src={getProxyUrl(data.user_info.face)} alt="" class="sc-face" />
    <div class="sc-info">
      <div class="sc-uname">{data.user_info.uname}</div>
      <div class="sc-price">¥{data.price}</div>
    </div>
  </div>
  <div class="sc-message">
    {data.message}
  </div>
  <div class="progress-bg">
    <div class="progress-fill" style="width: {remainingPercent}%"></div>
  </div>
</div>

<style>
  .sc-card {
    flex-shrink: 0;
    width: 220px; /* 固定宽度以便横向排列 */
    background: var(--overlay);
    border-radius: 6px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    border: 1px solid rgba(255, 255, 255, 0.1);
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.4);
    /* 强制所有文字为白色 */
    color: #ffffff !important; 
    animation: fadeIn 0.3s ease;
  }

  .sc-header {
    background: var(--sc-bg);
    padding: 8px;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .sc-face {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    border: 1.5px solid rgba(255, 255, 255, 0.5);
  }

  .sc-info {
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .sc-uname {
    font-size: 12px;
    font-weight: 700;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    color: #ffffff; /* 确保是白色 */
  }

  .sc-price {
    font-size: 11px;
    opacity: 0.9;
    font-weight: 600;
    color: #ffffff;
  }

  .sc-message {
    padding: 8px;
    font-size: 13px;
    line-height: 1.4;
    flex: 1;
    min-height: 40px;
    word-break: break-all;
    color: #ffffff; /* 确保是白色 */
  }

  .progress-bg {
    height: 3px;
    background: rgba(255, 255, 255, 0.1);
  }

  .progress-fill {
    height: 100%;
    background: #ffffff; /* 倒计时条也改为白色，更显眼 */
    transition: width 1s linear;
    box-shadow: 0 0 5px #ffffff;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(-10px); }
    to { opacity: 1; transform: translateY(0); }
  }
</style>
