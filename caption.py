#!/usr/bin/env python
# coding: utf-8

# Reference https://huggingface.co/microsoft/Florence-2-large/blob/main/sample_inference.ipynb

from transformers import AutoProcessor, AutoModelForCausalLM  
from PIL import Image
from io import BytesIO
import torch
from flask import Flask, request, jsonify

model_id = 'microsoft/Florence-2-large'
print("AutoModelForCausalLM.from_pretrained...")
model = AutoModelForCausalLM.from_pretrained(model_id, trust_remote_code=True, torch_dtype='auto').eval().cuda()
print("AutoProcessor.from_pretrained...")
processor = AutoProcessor.from_pretrained(model_id, trust_remote_code=True)

def run_example(task_prompt, image, text_input=None):
    if text_input is None:
        prompt = task_prompt
    else:
        prompt = task_prompt + text_input
    inputs = processor(text=prompt, images=image, return_tensors="pt").to('cuda', torch.float16)
    generated_ids = model.generate(
      input_ids=inputs["input_ids"].cuda(),
      pixel_values=inputs["pixel_values"].cuda(),
      max_new_tokens=1024,
      early_stopping=False,
      do_sample=False,
      num_beams=3,
    )
    generated_text = processor.batch_decode(generated_ids, skip_special_tokens=False)[0]
    parsed_answer = processor.post_process_generation(
        generated_text, 
        task=task_prompt, 
        image_size=(image.width, image.height)
    )

    return parsed_answer

app = Flask(__name__)

@app.route('/upload', methods=['POST'])
def upload_image():
    file = request.files['image']
    if file:
        
        # TODO Check size
        image_data = file.read()
        image = Image.open(BytesIO(image_data))

        print("run_example...")
        task_prompt = '<MORE_DETAILED_CAPTION>'
        results = run_example(task_prompt, image)

        return jsonify({'message': results})
    else:
        return jsonify({'message': 'No file uploaded'})

if __name__ == '__main__':
    app.run()