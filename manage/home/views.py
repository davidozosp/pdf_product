import re
import requests
from io import BytesIO
from django.shortcuts import render
from PyPDF2 import PdfReader, PdfWriter, PdfMerger
from django.conf import settings
from common.clean import clean_product_text
# Create your views here.


def home_view(request):
    
    if request.method == "POST":
        
        PORT_MANAGE = settings.PORT_MANAGE
        url_handler = f"http://manage:{PORT_MANAGE}/api/uploadfile"
        url_handler_split_sku = f"http://manage:{PORT_MANAGE}/api/split-sku"
        
        files = request.FILES
        fileupload = []
        for index, value in enumerate(files):
            fileupload.append(('files[]', (F"{index}.pdf", files[value].read(), 'application/pdf')))
            
            
        
        response_pdf2text = requests.request("POST", url_handler, files=fileupload)
        orders_json = response_pdf2text.json()
        orders = orders_json['orders']
        pdf_str = orders_json['pdf']
        pdf_bytes = pdf_str.encode()
        
        response_split_sku = requests.request("POST", url_handler_split_sku, json = {"orders": orders})
        
        
        return render(request, 'index.html', {})
    else:
        return render(request, 'index.html', {})