<script lang="ts">
    import SwitchButton from "./switchButton.svelte";
    import QuitButton from "./quitButton.svelte";
    import GithubButton from "./githubButton.svelte";
</script>

<div class="control-panel">
    <div class="control-section">
        <span>Start At Login</span>
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
    }

    .control-section {
        display: flex;
        align-items: center;
        justify-content: space-between;
    }

    .divider {
        height: 15px;
    }
</style>
