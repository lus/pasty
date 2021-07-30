// Encrypts a piece of text using AES-CBC and returns the HEX-encoded key, initialization vector and encrypted text
export async function encrypt(encryptionData, text) {
    const key = encryptionData.key;
    const iv = encryptionData.iv;

    const textBytes = aesjs.padding.pkcs7.pad(aesjs.utils.utf8.toBytes(text));

    const aes = new aesjs.ModeOfOperation.cbc(key, iv);
    const encrypted = aes.encrypt(textBytes);
    
    return {
        key: aesjs.utils.hex.fromBytes(key),
        iv: aesjs.utils.hex.fromBytes(iv),
        result: aesjs.utils.hex.fromBytes(encrypted)
    };
}

// Decrypts an encrypted piece of AES-CBC encrypted text
export async function decrypt(keyHex, ivHex, inputHex) {
    const key = aesjs.utils.hex.toBytes(keyHex);
    const iv = aesjs.utils.hex.toBytes(ivHex);
    const input = aesjs.utils.hex.toBytes(inputHex);

    const aes = new aesjs.ModeOfOperation.cbc(key, iv);
    const decrypted = aesjs.padding.pkcs7.strip(aes.decrypt(input));

    return aesjs.utils.utf8.fromBytes(decrypted);
}

// Creates encryption data from hex key and IV
export async function encryptionDataFromHex(keyHex, ivHex) {
    return {
        key: aesjs.utils.hex.toBytes(keyHex),
        iv: aesjs.utils.hex.toBytes(ivHex)
    };
}

// Generates encryption data to pass into the encrypt function
export async function generateEncryptionData() {
    return {
        key: await generateKey(),
        iv: generateIV()
    };
}

// Generates a new 256-bit AES-CBC key
async function generateKey() {
    const key = await crypto.subtle.generateKey({
        name: "AES-CBC",
        length: 256
    }, true, ["encrypt", "decrypt"]);

    const extracted = await crypto.subtle.exportKey("raw", key);
    return new Uint8Array(extracted);
}

// Generates a new cryptographically secure 16-byte array which is used as the initialization vector (IV) for AES-CBC
function generateIV() {
    return crypto.getRandomValues(new Uint8Array(16));
}