const HEX_REGEX = /^[0-9a-fA-F]+$/;

const normalizeHex = (value: string) =>
  value.startsWith('0x') || value.startsWith('0X') ? value.slice(2) : value;

export const bytesToHex = (bytes?: Uint8Array | null) => {
  if (!bytes || bytes.length === 0) return '0x';

  return `0x${Array.from(bytes)
    .map((byte) => byte.toString(16).padStart(2, '0'))
    .join('')}`;
};

export const hexToBytes = (value: string): Uint8Array => {
  const normalized = normalizeHex(value.trim());

  if (normalized.length === 0) {
    return new Uint8Array();
  }

  if (!HEX_REGEX.test(normalized) || normalized.length % 2 !== 0) {
    throw new Error(
      'Hex values must contain an even number of characters 0-9 or a-f'
    );
  }

  const bytes = new Uint8Array(normalized.length / 2);

  for (let i = 0; i < normalized.length; i += 2) {
    bytes[i / 2] = parseInt(normalized.slice(i, i + 2), 16);
  }

  return bytes;
};

export const parseByteInput = (value: string): Uint8Array => {
  const trimmed = value.trim();

  if (!trimmed) {
    return new Uint8Array();
  }

  try {
    return hexToBytes(trimmed);
  } catch (hexError) {
    try {
      return new TextEncoder().encode(trimmed);
    } catch {
      throw hexError instanceof Error
        ? hexError
        : new Error('Invalid byte input');
    }
  }
};
