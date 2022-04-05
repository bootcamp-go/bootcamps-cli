const libsodium = require('libsodium-wrappers');

// Compatible with the same `Uint8Array` arguments as `tweetsodium.seal()`
async function async_encrypt(messageBytes, publicKey) {
	await libsodium.ready;
	const d = libsodium.crypto_box_seal(messageBytes, publicKey);
	const e = Buffer.from(d).toString('base64');
	console.log(e);
}

// base64-encoded public key
const publicKey = process.argv.slice(2)[0];
const keyBytes = Buffer.from(publicKey, 'base64');

async_encrypt(Buffer.from(process.argv.slice(3)[0]), keyBytes);
