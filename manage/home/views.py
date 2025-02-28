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
        HOST_PARSER = settings.HOST_PARSER
        pages : list[tuple] = []
        merges = PdfMerger()
        for index, value in enumerate(request.FILES):
            reader = PdfReader(request.FILES[value])
            for _, page in enumerate(reader.pages):
                writer = PdfWriter()
                writer.add_page(page)
                pdf_buffer = BytesIO()
                writer.write(pdf_buffer)
                pages.append(('files[]', (F"{index}_{_}.pdf", pdf_buffer.getvalue(), 'application/pdf')))
                merges.append(pdf_buffer)
        merges.write("a.pdf")
        url = F"http://ocr-api:{HOST_PARSER}/pdf-parser"
        response = requests.request("GET", url, files=pages)
        responseText : list[dict] = response.json()
        
        orders :list = []
    
        for index_page, page in enumerate(responseText):
            contentOrder = page['ContentOrder']
            _product_lines = re.split(r'(?=\d+\.)', contentOrder)
            product_lines = [pro.strip() for pro in _product_lines if pro != ""]
            for index_product, product in enumerate(product_lines):
                content_order, sl = clean_product_text(product)
                orders.append({
                    "page" : index_page,
                    "order_id" : page["OrderId"],
                    "content": content_order,
                    "quantity": sl,
                    "total_product": len(product_lines),
                })
                
        for order in orders:
            print(order)
        
        
        return render(request, 'index.html', {})
    else:
        return render(request, 'index.html', {})