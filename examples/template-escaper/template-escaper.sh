# remove comments
# add escape function to pipeline
EOF State
<p><!-- HTML comment -->abc</p>
.Content
EOF State
<p>abc</p>.Content | EscapeHtml

Program exited.
