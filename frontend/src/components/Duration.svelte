<script lang="ts">
  import { searchListstore } from "../stores";
  const { hours, minutes, seconds, setHours, setMinutes, setSeconds } =
    searchListstore;

  function handleHoursChange(event) {
    setHours(Math.min(2, Math.max(0, event.target.value)));
  }

  function handleMinutesChange(event) {
    setMinutes(Math.min(59, Math.max(0, event.target.value)));
  }

  function handleSecondsChange(event) {
    setSeconds(Math.min(59, Math.max(0, event.target.value)));
  }
</script>

<div class="max-w-sm mx-auto p-4 bg-obsternary rounded-md">
  <h2 class="text-xl font-semibold mb-4 text-white">Duration</h2>

  <div class="flex space-x-4 mb-4">
    <div class="flex-1">
      <label for="hours" class="block text-sm font-medium text-white"
        >Hours</label
      >
      <input
        type="number"
        id="hours"
        bind:value={$hours}
        on:input={handleHoursChange}
        min="0"
        max="2"
        class="mt-1 block w-full px-3 py-2 border border-white rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 text-black"
      />
    </div>

    <div class="flex-1">
      <label for="minutes" class="block text-sm font-medium text-white"
        >Minutes</label
      >
      <input
        type="number"
        id="minutes"
        bind:value={$minutes}
        on:input={handleMinutesChange}
        min="0"
        max="59"
        class="mt-1 block w-full px-3 py-2 border borde-white rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 text-black"
      />
    </div>

    <div class="flex-1">
      <label for="seconds" class="block text-sm font-medium text-white"
        >Seconds</label
      >
      <input
        type="number"
        id="seconds"
        bind:value={$seconds}
        on:input={handleSecondsChange}
        min="0"
        max="59"
        class="mt-1 block w-full px-3 py-2 border border-white rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 text-black"
      />
    </div>
  </div>

  <p class="text-white">
    <span class="font-medium">{$hours ? $hours : 0} hours</span>,
    <span class="font-medium">{$minutes ? $minutes : 0} minutes</span>,
    <span class="font-medium">{$seconds ? $seconds : 0} seconds</span>
  </p>
</div>
