/**
 * Linearly interpolates between two values.
 * @param {number} a The start value.
 * @param {number} b The end value.
 * @param {number} t The interpolation factor (0-1).
 * @returns {number} The interpolated value.
 */
export function lerp(a, b, t) {
    return a + (b - a) * t;
}

/**
 * Linearly interpolates between two angles.
 * @param {number} a The start angle (in radians).
 * @param {number} b The end angle (in radians).
 * @param {number} t The interpolation factor (0-1).
 * @returns {number} The interpolated angle (in radians).
 */
export function lerpAngle(a, b, t) {
    return Math.atan2(
        (1 - t) * Math.sin(a) + t * Math.sin(b),
        (1 - t) * Math.cos(a) + t * Math.cos(b)
    );
}

/**
 * Formats a number of bits into a human-readable string.
 * @param {number} numBits The number of bits.
 * @returns {string} The formatted string.
 */
export function formatBits(numBits) {
    if (numBits < 1024) {
        return `${numBits.toFixed(2)} b`;
    } else if (numBits < 1024 * 1024) {
        return `${(numBits / 1024).toFixed(2)} Kb`;
    } else if (numBits < 1024 * 1024 * 1024) {
        return `${(numBits / (1024 * 1024)).toFixed(2)} Mb`;
    } else {
        return `${(numBits / (1024 * 1024 * 1024)).toFixed(2)} Gb`;
    }
}