<script lang="ts">
    import { exit } from "@tauri-apps/plugin-process";

    const quitApp = async () => {
        await exit(0);
    };
</script>

<button class="quit-button" onclick={quitApp} aria-label="Quit">Quit</button>

<style>
    :global(button) {
        background: none;
        color: inherit;
        border: none;
        padding: 0;
        font: inherit;
        cursor: pointer;
        outline: inherit;
    }

    .quit-button {
        background-color: red;
        padding-top: 5px;
        padding-bottom: 5px;
        width: 48px;
        border-radius: 16px;
        transition: background-color 100ms;
    }

    .quit-button:hover {
        background-color: darkred;
    }
</style>
