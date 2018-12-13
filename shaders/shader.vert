#version 410

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec2 vertTexCoord;
in vec3 vertNormal;

out vec2 fragTexCoord;
out vec3 fragNormal;
out vec3 fragPosition;

void main() {
  mat4 modelview = camera * model;

  gl_Position = projection * modelview * vec4(vert, 1);
  fragTexCoord = vertTexCoord;
  fragPosition = vert;
  fragNormal = vertNormal;
}
