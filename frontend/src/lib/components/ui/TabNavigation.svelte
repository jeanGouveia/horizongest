<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let tabs: { id: string; label: string }[] = [];
  export let activeTab: string = tabs[0]?.id ?? '';
  export let size: 'sm' | 'md' | 'lg' = 'md';

  const dispatch = createEventDispatcher();

  function handleTabClick(tabId: string) {
    activeTab = tabId;
    dispatch('tabChange', { tabId });
  }

  $: activeIndex = tabs.findIndex(t => t.id === activeTab);
</script>

<div class="tab-navigation" class:size>
  <div class="tab-list" role="tablist">
    {#each tabs as tab (tab.id)}
      <button
        class="tab-button"
        class:active={activeTab === tab.id}
        role="tab"
        aria-selected={activeTab === tab.id}
        aria-controls={`panel-${tab.id}`}
        id={`tab-${tab.id}`}
        onclick={() => handleTabClick(tab.id)}
      >
        {tab.label}
      </button>
    {/each}
  </div>
  <div class="tab-indicator" style="transform: translateX({activeIndex * 100}%)" />
</div>

<style>
  .tab-navigation {
    position: relative;
    border-bottom: 1px solid #f1f5f9;
  }

  .tab-list {
    display: flex;
    position: relative;
    background: #ffffff;
  }

  .tab-button {
    position: relative;
    padding: 1rem 1.5rem;
    background: transparent;
    border: none;
    border-bottom: 2px solid transparent;
    color: #64748b;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: color 0.15s ease;
    white-space: nowrap;
  }

  .tab-button:hover {
    color: #0f172a;
  }

  .tab-button.active {
    color: #6366f1;
    border-bottom-color: #6366f1;
  }

  .tab-button:focus {
    outline: none;
    box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.2);
  }

  .tab-indicator {
    position: absolute;
    bottom: -1px;
    left: 0;
    height: 2px;
    background: #6366f1;
    transition: transform 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    pointer-events: none;
  }

  .tab-navigation.sm .tab-button {
    padding: 0.75rem 1rem;
    font-size: 0.8125rem;
  }

  .tab-navigation.lg .tab-button {
    padding: 1.25rem 2rem;
    font-size: 1rem;
  }
</style>
