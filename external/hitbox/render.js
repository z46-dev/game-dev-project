import { lerp } from "./utils.js";

/**
 * Gets the UI scale for a given canvas.
 * @param {HTMLCanvasElement|OffscreenCanvas} c The canvas to get the UI scale for.
 * @returns {number} The UI scale.
 */
export function uiScale(c) {
    if (c.height > c.width) {
        return c.height / 1080;
    }

    return c.width / 1920;
}

/**
 * Gets the game scale for a given canvas and field of view.
 * @param {HTMLCanvasElement|OffscreenCanvas} c The canvas to get the game scale for.
 * @param {number} fov The field of view.
 * @returns {number} The game scale.
 */
export function gameScale(c, fov) {
    return Math.max(c.width / fov, c.height / fov / 1080 * 1920);
}

const colorCache = new Map();

/**
 * Mixes two colors together.
 * @param {string} primary The primary color in hex format.
 * @param {string} secondary The secondary color in hex format.
 * @param {number} [amount=0.5] The amount of the secondary color to mix in (0-1).
 * @returns {string} The mixed color in hex format.
 */
export function mixColors(primary, secondary, amount = .5) {
    const key = `${primary}${secondary}${amount}`;

    if (colorCache.has(key)) {
        return colorCache.get(key);
    }

    const pr = parseInt(primary.slice(1), 16);
    const sc = parseInt(secondary.slice(1), 16);

    const hex = `#${(
        1 << 24 |
        (lerp((pr >> 16) & 255, (sc >> 16) & 255, amount) | 0) << 16 |
        (lerp((pr >> 8) & 255, (sc >> 8) & 255, amount) | 0) << 8 |
        (lerp(pr & 255, sc & 255, amount) | 0)
    ).toString(16).slice(1)}`;

    colorCache.set(key, hex);
    return hex;
}

/**
 * Renders text on a given canvas.
 * @param {CanvasRenderingContext2D} c The canvas context to render the text on.
 * @param {string} text The text to render.
 * @param {number} x The x position to render the text at.
 * @param {number} y The y position to render the text at.
 * @param {number} size The font size of the text.
 * @param {string} [fill="#FFFFFF"] The fill color of the text in hex format.
 * @returns {number} The width of the rendered text.
 */
export function text(c, text, x, y, size, fill = "#FFFFFF") {
    c.fillStyle = fill;
    c.strokeStyle = mixColors(fill, "#000000", .2);
    c.lineWidth = size * .15;
    c.font = `bold ${size}px sans-serif`;

    c.strokeText(text, x, y);
    c.fillText(text, x, y);

    return c.measureText(text).width;
}

/**
 * Sets the fill and stroke style for a given canvas context.
 * @param {CanvasRenderingContext2D} c The canvas context to set the style for.
 * @param {string} fillColor The fill color in hex format.
 * @param {number} [lineWidth=0.1] The line width for the stroke.
 * @param {number} [darkening=0.2] The amount to darken the stroke color (0-1).
 */
export function setStyle(c, fillColor, lineWidth = .1, darkening = .2) {
    c.fillStyle = fillColor;
    c.strokeStyle = mixColors(fillColor, "#000000", darkening);
    c.lineWidth = lineWidth;
}

/**
 * Sets the fill and stroke color for a given canvas context.
 * @param {CanvasRenderingContext2D} c The canvas context to set the color for.
 * @param {string} fillColor The fill color in hex format.
 * @param {number} [darkening=0.2] The amount to darken the stroke color (0-1).
 */
export function setColor(c, fillColor, darkening = .2) {
    c.fillStyle = fillColor;
    c.strokeStyle = mixColors(fillColor, "#000000", darkening);
}

/**
 * Draws a rounded rectangle on a given canvas context.
 * @param {CanvasRenderingContext2D} c The canvas context to draw the rounded rectangle on.
 * @param {number} cx The center x position of the rectangle.
 * @param {number} cy The center y position of the rectangle.
 * @param {number} width The width of the rectangle.
 * @param {number} height The height of the rectangle.
 * @param {number} radius The corner radius of the rectangle.
 */
export function roundedRectangle(c, cx, cy, width, height, radius) {
    c.moveTo(cx - width / 2 + radius, cy - height / 2);
    c.roundRect(cx - width / 2, cy - height / 2, width, height, radius);
}

/**
 * Draws a horizontal bar on a given canvas context.
 * @param {CanvasRenderingContext2D} c The canvas context to draw the bar on.
 * @param {number} x1 The starting x position of the bar.
 * @param {number} x2 The ending x position of the bar.
 * @param {number} y The y position of the bar.
 * @param {number} thickness The thickness of the bar.
 * @param {string} [color="#555555"] The color of the bar in hex format.
 */
export function drawBar(c, x1, x2, y, thickness, color = "#555555") {
    c.strokeStyle = color;
    c.lineWidth = thickness;

    c.beginPath();
    c.moveTo(x1, y);
    c.lineTo(x2, y);
    c.closePath();
    c.stroke();
}

/**
 * Draws a progress bar on a given canvas context.
 * @param {CanvasRenderingContext2D} c The canvas context to draw the progress bar on.
 * @param {number} centerX The center x position of the progress bar.
 * @param {number} y The y position of the progress bar.
 * @param {number} width The width of the progress bar.
 * @param {number} thickness The thickness of the progress bar.
 * @param {number} borderWidth The border width of the progress bar.
 * @param {number} progress The progress value (0-1).
 * @param {string} [fillColor="#A5A5A5"] The fill color of the progress bar in hex format.
 * @param {string} [backColor="#555555"] The background color of the progress bar in hex format.
 */
export function drawProgressBar(c, centerX, y, width, thickness, borderWidth, progress, fillColor = "#A5A5A5", backColor = "#555555") {
    const oldLineCap = c.lineJoin;
    c.lineJoin = "round";

    // Background
    c.strokeStyle = backColor;
    c.lineWidth = thickness + borderWidth;
    c.beginPath();
    c.moveTo(centerX - width / 2, y);
    c.lineTo(centerX + width / 2, y);
    c.closePath();
    c.stroke();

    // Progress
    c.strokeStyle = fillColor;
    c.lineWidth = thickness;
    c.beginPath();
    c.moveTo(centerX - width / 2, y);
    c.lineTo(centerX - width / 2 + width * progress, y);
    c.closePath();
    c.stroke();

    c.lineJoin = oldLineCap;
}

/**
 * So Canvas is very slow once you apply save, restore, translate, scale, rotate.
 * To work around this, we can use ctx.setTransform for things we can't calculate ourselves.
 * For simple lines, we can calculate this. So these functions are huge wins in performance.
 */

/**
 * Applies a transformation to a given canvas context.
 * @param {CanvasRenderingContext2D} c The canvas context to apply the transformation to.
 * @param {number} x The x position to translate to.
 * @param {number} y The y position to translate to.
 * @param {number} size The scale factor to apply.
 * @param {number} angle The rotation angle in radians.
 * @returns {DOMMatrix} The previous transformation matrix.
 */
export function transform(c, x, y, size, angle) {
    const old = c.getTransform();
    const cs = Math.cos(angle) * size;
    const sn = Math.sin(angle) * size;
    c.setTransform(cs, sn, -sn, cs, x, y);
    return old;
}

/**
 * Restores the default transformation on a given canvas context.
 * @param {CanvasRenderingContext2D} c The canvas context to restore the transformation to.
 */
export function resetTransform(c) {
    c.setTransform(1, 0, 0, 1, 0, 0);
}

/**
 * Applies a non-rotational transformation to a given canvas context.
 * @param {CanvasRenderingContext2D} c The canvas context to apply the transformation to.
 * @param {number} x The x position to translate to.
 * @param {number} y The y position to translate to.
 * @param {number} scale The scale factor to apply.
 * @returns {DOMMatrix} The previous transformation matrix.
 */
export function transformNoRotate(c, x, y, scale) {
    const old = c.getTransform();
    c.setTransform(scale, 0, 0, scale, x, y);
    return old;
}

/**
 * Optimized moveTo function to avoid using save/restore/translate/rotate.
 * @param {CanvasRenderingContext2D} c The canvas context to move the path to.
 * @param {number} x The x position to move to.
 * @param {number} y The y position to move to.
 * @param {number} cs The cosine of the rotation angle.
 * @param {number} sn The sine of the rotation angle.
 * @param {number} [ox=0] The origin x offset.
 * @param {number} [oy=0] The origin y offset.
 */
export function moveTo(c, x, y, cs, sn, ox = 0, oy = 0) {
    c.moveTo(ox + x * cs - y * sn, oy + x * sn + y * cs);
}

/**
 * Optimized lineTo function to avoid using save/restore/translate/rotate.
 * @param {CanvasRenderingContext2D} c The canvas context to draw the line to.
 * @param {number} x The x position to draw the line to.
 * @param {number} y The y position to draw the line to.
 * @param {number} cs The cosine of the rotation angle.
 * @param {number} sn The sine of the rotation angle.
 * @param {number} [ox=0] The origin x offset.
 * @param {number} [oy=0] The origin y offset.
 */
export function lineTo(c, x, y, cs, sn, ox = 0, oy = 0) {
    c.lineTo(ox + x * cs - y * sn, oy + x * sn + y * cs);
}

/**
 * Draws a circle with a border on a given canvas context.
 * @param {CanvasRenderingContext2D} c The canvas context to draw the circle on.
 * @param {number} x The x position of the circle center.
 * @param {number} y The y position of the circle center.
 * @param {number} radius The radius of the circle.
 * @param {number} borderWidth The width of the border.
 */
export function circleWithBorder(c, x, y, radius, borderWidth) {
    const fill = c.fillStyle;

    c.fillStyle = c.strokeStyle;
    c.beginPath();
    c.arc(x, y, radius + borderWidth, 0, Math.PI * 2);
    c.fill();

    c.fillStyle = fill;
    c.beginPath();
    c.arc(x, y, radius, 0, Math.PI * 2);
    c.fill();
}