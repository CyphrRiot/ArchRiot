update-desktop-database ~/.local/share/applications

# Configure default applications for file types
echo "📄 Configuring default applications for file types..."

# Set image viewer (imv) for all image formats
for mimetype in image/png image/jpeg image/gif image/webp image/bmp image/tiff; do
    xdg-mime default imv.desktop "$mimetype"
done

# Set PDF viewer
xdg-mime default org.gnome.Papers.desktop application/pdf

# Set default browser
xdg-settings set default-web-browser brave-browser.desktop
for scheme in x-scheme-handler/http x-scheme-handler/https; do
    xdg-mime default brave-browser.desktop "$scheme"
done

# Set text editor for text and markdown
for mimetype in text/plain text/markdown text/x-markdown application/x-markdown; do
    xdg-mime default org.gnome.TextEditor.desktop "$mimetype"
done

# Set video player (mpv) for all video formats
for mimetype in video/mp4 video/x-msvideo video/x-matroska video/x-flv video/x-ms-wmv video/mpeg video/ogg video/webm video/quicktime video/3gpp video/3gpp2 video/x-ms-asf video/x-ogm+ogg video/x-theora+ogg application/ogg; do
    xdg-mime default mpv.desktop "$mimetype"
done

echo "✓ Default applications configured"
