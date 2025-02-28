import re

def clean_product_text(text):
    """
    - Xoá số thứ tự đầu dòng (1., 2., 3., ...).
    - Tách giá trị số cuối cùng sang cột SL.
    - Xóa tất cả ký tự từ dấu ',' cuối cùng đến ngay trước số đó.
    - Định dạng SL là số (int).
    """
    # Xóa số thứ tự đầu dòng (nếu có)
    text = re.sub(r"^\d+\.\s*", "", text)  # Xóa "1. ", "2. ", "3. ", ...

    sl_value = None
    match = re.search(r",\s*([^,]*?)(\d+)\s*$", text)

    if match:
        sl_value = int(match.group(2))  # Lấy số cuối cùng làm SL
        text = text[: match.start()].strip()  # Cắt bỏ nội dung từ dấu ',' cuối cùng

    return text.strip(), sl_value