#version 410

uniform sampler2D tex;
in vec2 fragTexCoord;
out vec4 outputColor;

void main() {
    outputColor = texture(tex, vec2(fragTexCoord.x, 1.0-fragTexCoord.y));
}