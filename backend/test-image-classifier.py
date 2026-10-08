#!/usr/bin/env python3
"""Repeat an image-classifier request to compare sampling stability.

Example:
  python3 backend/test-image-classifier.py --runs 5 --temperature 0.2
  python3 backend/test-image-classifier.py --runs 5 --temperature 0 --seed 42
"""

import argparse
import base64
import json
import sys
import time
import urllib.request
from pathlib import Path

PROMPT = """Return only one valid JSON object, with exactly these fields: {"description":"one concise sentence","hashtags":["short-tag"],"text_tags":["visible-word"]}.
Describe what is visibly happening in the image. Hashtags must be 3-6 short lowercase visual or meme-context tags, without # or explanations. text_tags must contain only complete, clearly readable words or short phrases visible in the image, kept as close to the exact spelling as possible, lowercase, without punctuation. Never abbreviate, truncate, or guess a partial word; omit anything uncertain. Use an empty array only when no text is visible.
Always inspect the image for text_tags."""


def post(url: str, payload: dict) -> dict:
    body = json.dumps(payload).encode()
    request = urllib.request.Request(
        url,
        data=body,
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    with urllib.request.urlopen(request, timeout=120) as response:
        return json.load(response)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("image", type=Path)
    parser.add_argument("--base-url", default="https://image-classifier.dev.initialed85.cc/v1")
    parser.add_argument("--runs", type=int, default=5)
    parser.add_argument("--temperature", type=float, default=0.2)
    parser.add_argument("--seed", type=int)
    parser.add_argument("--top-p", type=float)
    parser.add_argument("--top-k", type=int)
    parser.add_argument("--max-tokens", type=int, default=600)
    args = parser.parse_args()

    image_type = "image/png"
    image = base64.b64encode(args.image.read_bytes()).decode()
    model_response = urllib.request.urlopen(args.base_url.rstrip("/") + "/models", timeout=30)
    model = json.load(model_response)["data"][0]["id"]

    settings = {
        "temperature": args.temperature,
        "max_tokens": args.max_tokens,
        "response_format": {"type": "json_object"},
        "chat_template_kwargs": {"enable_thinking": False},
    }
    if args.seed is not None:
        settings["seed"] = args.seed
    if args.top_p is not None:
        settings["top_p"] = args.top_p
    if args.top_k is not None:
        settings["top_k"] = args.top_k

    payload = {
        "model": model,
        "messages": [
            {"role": "system", "content": "You analyze images for a private meme library. Never mention that you are an AI and never wrap JSON in markdown."},
            {"role": "user", "content": [
                {"type": "text", "text": PROMPT},
                {"type": "image_url", "image_url": {"url": f"data:{image_type};base64,{image}"}},
            ]},
        ],
        **settings,
    }

    print(json.dumps({"model": model, "settings": settings, "runs": args.runs}, indent=2))
    outputs = []
    for index in range(args.runs):
        started = time.monotonic()
        response = post(args.base_url.rstrip("/") + "/chat/completions", payload)
        content = response["choices"][0]["message"].get("content", "")
        outputs.append(content)
        print(f"\n--- run {index + 1} ({time.monotonic() - started:.2f}s) ---\n{content}")

    print(f"\nidentical outputs: {len(set(outputs)) == 1} ({len(set(outputs))} unique)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
