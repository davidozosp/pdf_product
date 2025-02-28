from PyPDF2 import PdfReader
from PyPDF2 import PdfWriter


def split_pdf(pdf_path):
    reader = PdfReader(pdf_path)
    n = len(reader.pages)
    for i in range(n):
        writer = PdfWriter()
        page = reader.pages[i]
        writer.add_page(page)
        with open(f"input/page_{i}.pdf", "wb") as f:
            writer.write(f)

if __name__ == '__main__':
    pdf_path = 'merged_output.pdf'
    split_pdf(pdf_path)