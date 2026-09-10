import tkinter as tk
from tkinter import filedialog, messagebox
from PIL import Image, ImageTk


class ImageInspectorApp:

    def __init__(self, root):
        self.root = root
        self.root.title("Image Inspector - Startup")
        self.center_window(self.root, 600, 400)

        # Startup label
        tk.Label(root,
            text="Image\nOpen an image to begin.",
            font=("Arial", 18)
        ).pack(pady=40)

        # Open image button
        tk.Button(root,
            text="Open Image",
            font=("Arial", 14),
            width=15,
            command=self.open_image
        ).pack(pady=20)

        self.image_path = None
        self.image = None


    # Center a window on the screen

    def center_window(self, window, width, height):

        # screen dimensions
        screen_width = window.winfo_screenwidth()
        screen_height = window.winfo_screenheight()

        # Calculate centered position
        x = (screen_width - width) // 2
        y = (screen_height - height) // 2

        # Set size
        window.geometry(f"{width}x{height}+{x}+{y}")

  
    # Open image
  
    def open_image(self):

        file_path = filedialog.askopenfilename(title="Select an Image",
            filetypes=[("Image Files", "*.png *.jpg *.jpeg *.bmp *.gif")])

        if not file_path:
            return

        self.image_path = file_path

        try:
            self.image = Image.open(file_path)

        except Exception as e:
            messagebox.showerror("Error",f"Could not open image: {e}")
            return

        # Open inspector window
        self.open_inspector_window()


    # Inspector window
    def open_inspector_window(self):

        inspector = tk.Toplevel(self.root)
        inspector.title("Image Inspector")

        # Larger + centered
        self.center_window(inspector, 700, 650)

        # Resize image for preview
        preview_image = self.image.copy()
        preview_image.thumbnail((500, 400))

        img_preview = ImageTk.PhotoImage(preview_image)

        # Show image
        tk.Label(inspector, image=img_preview).pack(pady=20)

        # Keep reference alive
        inspector.img_preview = img_preview

        # Buttons
        tk.Button(inspector, text="Show Info",
            width=20,
            command=lambda: self.show_info(inspector)
        ).pack(pady=5)

        tk.Button(inspector, text="Convert to Grayscale",
            width=20,
            command=lambda: self.convert_grayscale(inspector)
        ).pack(pady=5)

        tk.Button(inspector, text="Resize Image",
            width=20,
            command=lambda: self.resize_image(inspector)
        ).pack(pady=5)

    # Show image information

    def show_info(self, parent):

        info_win = tk.Toplevel(parent)
        info_win.title("Image Info")
        self.center_window(info_win, 600, 300)

        w, h = self.image.size

        tk.Label(info_win, text=f"Image Path:\n{self.image_path}",
            wraplength=550,
            justify="left"
        ).pack(pady=15)

        tk.Label(info_win, text=f"Dimensions: {w} x {h}",
            font=("Arial", 12)
        ).pack(pady=5)

        tk.Label(info_win, text=f"Format: {self.image.format}",
            font=("Arial", 12)
        ).pack(pady=5)

  
    def convert_grayscale(self, parent):

        gray_win = tk.Toplevel(parent)
        gray_win.title("Grayscale Image")
        self.center_window(gray_win, 600, 600)
        gray_img = self.image.convert("L")
        gray_img.thumbnail((500, 450))
        gray_preview = ImageTk.PhotoImage(gray_img)

        tk.Label(gray_win, image=gray_preview).pack(pady=20)

        # Keep reference alive
        gray_win.gray_preview = gray_preview

  
    # Resize image

    def resize_image(self, parent):

        resize_win = tk.Toplevel(parent)
        resize_win.title("Resize Image")
        self.center_window(resize_win, 500, 600)

        tk.Label(resize_win,text="Enter new width:").pack(pady=(20, 5))

        width_entry = tk.Entry(resize_win)
        width_entry.pack()

        tk.Label(resize_win, text="Enter new height:").pack(pady=(15, 5))

        height_entry = tk.Entry(resize_win)
        height_entry.pack()

        preview_label = tk.Label(resize_win)
        preview_label.pack(pady=20)

        def apply_resize():

            try:
                new_w = int(width_entry.get())
                new_h = int(height_entry.get())

                if new_w <= 0 or new_h <= 0:
                    raise ValueError("Width and height must be greater than zero.")

                resized = self.image.resize((new_w, new_h))

                # Prevent extremely large previews
                preview_image = resized.copy()
                preview_image.thumbnail((400, 350))
                preview = ImageTk.PhotoImage(preview_image)
                preview_label.config(image=preview)

                # Keep reference alive
                preview_label.preview = preview

            except Exception as e:

                messagebox.showerror("Error", f"Resize failed: {e}")

        tk.Button(resize_win, text="Apply",
            width=15,
            command=apply_resize
        ).pack(pady=10)



if __name__ == "__main__":

    root = tk.Tk()
    app = ImageInspectorApp(root)

    root.mainloop()
