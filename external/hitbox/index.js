import * as renderLib from "./render.js";

/** @type {HTMLCanvasElement} */
const canvas = document.querySelector("canvas#gameCanvas");
const ctx = canvas.getContext("2d", {
    alpha: false,
    desynchronized: true
});

function resize() {
    canvas.width = window.innerWidth * devicePixelRatio;
    canvas.height = window.innerHeight * devicePixelRatio;

    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
}

resize();
window.addEventListener("resize", resize);

let objectX = 0,
    objectY = 0,
    objectSize = 500,
    mouseX = 0,
    mouseY = 0,
    relativeMouseX = 0,
    relativeMouseY = 0,
    points = [];

canvas.addEventListener("mousemove", (e) => {
    const rect = canvas.getBoundingClientRect();
    mouseX = (e.clientX - rect.left) * devicePixelRatio;
    mouseY = (e.clientY - rect.top) * devicePixelRatio;

    relativeMouseX = +((mouseX - objectX) / objectSize).toFixed(3);
    relativeMouseY = +((mouseY - objectY) / objectSize).toFixed(3);
});

canvas.addEventListener("mousedown", (e) => {
    switch (e.button) {
        case 0: // Left click (add point)
            points.push({
                x: relativeMouseX,
                y: relativeMouseY
            });
            console.log("Added point:", relativeMouseX, relativeMouseY);
            break;
        case 2: // Right click (remove closest within 25px)
            let closestIndex = -1;
            let closestDistance = 25 / objectSize; // 25 pixels in object space

            for (let i = 0; i < points.length; i++) {
                const point = points[i];
                const dx = point.x - relativeMouseX;
                const dy = point.y - relativeMouseY;
                const dist = Math.sqrt(dx * dx + dy * dy);

                if (dist < closestDistance) {
                    closestDistance = dist;
                    closestIndex = i;
                }
            }

            if (closestIndex !== -1) {
                const removedPoint = points.splice(closestIndex, 1)[0];
                console.log("Removed point:", removedPoint.x, removedPoint.y);
            }
            break;
    }
});

function loadImage(src) {
    return new Promise((resolve, reject) => {
        const img = new Image();
        img.onload = () => resolve(img);
        img.onerror = reject;
        img.src = src;
    });
}

const img = await loadImage("/assets/ships/parseval.png");
console.log("Loaded image:", img.width, "x", img.height);

function draw() {
    requestAnimationFrame(draw);

    ctx.clearRect(0, 0, canvas.width, canvas.height);

    objectX = canvas.width / 2;
    objectY = canvas.height / 2;

    renderLib.setStyle(ctx, "#FFFFFF", 2 / objectSize);
    renderLib.transform(ctx, objectX, objectY, objectSize, 0);

    ctx.drawImage(img, -1, -1, 2, 2);

    renderLib.resetTransform(ctx);

    renderLib.setStyle(ctx, "#FF0000", 0.1);
    renderLib.transformNoRotate(ctx, mouseX, mouseY, 25); // crosshair
    ctx.beginPath();
    ctx.moveTo(-1, 0);
    ctx.lineTo(1, 0);
    ctx.moveTo(0, -1);
    ctx.lineTo(0, 1);
    ctx.closePath();
    ctx.stroke();
    renderLib.resetTransform(ctx);

    const path = new Path2D();
    renderLib.setColor(ctx, "#00FF00");
    let i = 0;
    for (const point of points) {
        renderLib.circleWithBorder(ctx, objectX + point.x * objectSize, objectY + point.y * objectSize, 5, 2);
        if (i++ === 0) {
            path.moveTo(objectX + point.x * objectSize, objectY + point.y * objectSize);
        } else {
            path.lineTo(objectX + point.x * objectSize, objectY + point.y * objectSize);
        }
    }

    renderLib.setStyle(ctx, "#00FF00", 2);
    ctx.globalAlpha = 0.25;
    ctx.fill(path);
    ctx.globalAlpha = 1.0;
    ctx.stroke(path);

    renderLib.text(ctx, `Mouse: (${relativeMouseX}, ${relativeMouseY})`, canvas.width / 2, 30, 24, "#FFFFFF");
}

console.log("Starting hitbox render");
draw();

// Print out []*util.Vector2D{util.Vector(x, y), ...}
window.exportHitbox = function() {
    let output = "[]*util.Vector2D{\n";
    for (const point of points) {
        output += `    util.Vector(${point.x}, ${point.y}),\n`;
    }
    output += "}";
    console.log(output);
}