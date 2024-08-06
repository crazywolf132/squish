const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const fetch = require('node-fetch');

const version = require('./package.json').version;
const binPath = path.join(__dirname, 'bin');

async function getLatestRelease() {
    const response = await fetch('https://api.github.com/repos/foxycorps/squish/releases/latest');
    const data = await response.json();
    return data.tag_name.replace('v', '');
}

async function install() {
    const latestVersion = await getLatestRelease();
    const platform = process.platform;
    const arch = process.arch === 'x64' ? 'amd64' : process.arch;

    let filename;
    switch (platform) {
        case 'linux':
            filename = `squish-linux-${arch}`;
            break;
        case 'darwin':
            filename = `squish-darwin-${arch}`;
            break;
        case 'win32':
            filename = `squish-windows-${arch}.exe`;
            break;
        default:
            throw new Error(`Unsupported platform: ${platform}`);
    }

    const url = `https://github.com/foxycorps/squish/releases/download/v${latestVersion}/${filename}`;
    const outputPath = path.join(binPath, platform === 'win32' ? 'squish.exe' : 'squish');

    if (!fs.existsSync(binPath)) {
        fs.mkdirSync(binPath, { recursive: true });
    }

    console.log(`Downloading Squish v${latestVersion} for ${platform} ${arch}...`);

    try {
        execSync(`curl -L ${url} -o ${outputPath}`);
        fs.chmodSync(outputPath, 0o755); // Make the file executable
        console.log('Squish has been installed successfully!');
    } catch (error) {
        console.error('Failed to download Squish:', error);
        process.exit(1);
    }
}

install().catch(console.error);