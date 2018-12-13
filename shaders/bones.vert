#version 410

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
uniform mat4 bones[30];

in vec3 vert;
in vec2 vertTexCoord;
in vec3 vertNormal;
in ivec4 vertBones;
in vec4 vertWeight;

out vec2 fragTexCoord;
out vec3 fragNormal;

void main() {

  mat4 BoneTransform = bones[vertBones[0]] * vertWeight[0];
  BoneTransform += bones[vertBones[1]] * vertWeight[1];
  BoneTransform += bones[vertBones[2]] * vertWeight[2];
  BoneTransform += bones[vertBones[3]] * vertWeight[3];

  fragNormal = (BoneTransform * vec4(vertNormal, 0.0)).xyz;
  // fragNormal = vertNormal;

  fragTexCoord = vertTexCoord;
  
  gl_Position = projection * camera * model * BoneTransform * vec4(vert, 1);
}
