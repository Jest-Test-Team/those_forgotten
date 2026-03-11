import { copyFileSync, mkdirSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(fileURLToPath(import.meta.url));
const source = resolve(root, "../openapi.yaml");
const destination = resolve(root, "../dist/openapi.yaml");

mkdirSync(resolve(root, "../dist"), { recursive: true });
copyFileSync(source, destination);
