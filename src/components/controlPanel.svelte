<script lang="ts">
    import SwitchButton from "./switchButton.svelte";
    import QuitButton from "./quitButton.svelte";
    import GithubButton from "./githubButton.svelte";
</script>

<div class="control-panel">
    <div class="input-field">
        <span class="input-text">Start At Login</span>
        <SwitchButton />
    </div>
    <div class="divider"></div>
    <div class="control-section bottom-bar">
        <GithubButton />
        <QuitButton />
    </div>
</div>

<style>
    .control-panel {
        display: flex;
        flex-direction: column;
        font-size: 0.85rem;
        position: relative;
        height: 100%;
    }
    
    .input-field {
        display: flex;
        flex-direction: column;
    }

    .input-text {
        font-size: 16px;
        font-weight: 600;
        margin-bottom: 6px;
    }

    .control-section {
        display: flex;
        align-items: center;
        justify-content: space-between;
    }

    .divider {
        flex-grow: 1;
    }
</style>
