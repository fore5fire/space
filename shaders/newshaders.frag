#version 330 core
// #version 410

uniform sampler2D tex;
uniform mat4 camera;

in vec2 fragTexCoord;
in vec3 fragNormal;

out vec4 outputColor;

void main() {

  vec3 fragCamera = vec3(camera[3][0], camera[3][1], camera[3][2]);
  
  vec3 lightPosition = vec3(0,1,0);
  vec3 diffuseLightColor = vec3(1,1,1);
  vec3 specularObjectColor = vec3(1,1,1);
  vec3 ambientLight = vec3(0.2,0.2,0.2);
  float phongExponent = 10;

  vec3 color = texture(tex, vec2(fragTexCoord.x, 1.0-fragTexCoord.y)).rgb;
  vec3 fragPosition = vec3(gl_FragCoord.x, gl_FragCoord.y, gl_FragCoord.z);
  vec3 cameraDirection = normalize(fragCamera - fragPosition);
  vec3 lightDirection = normalize(lightPosition - fragPosition);

  vec3 ca = color * ambientLight;
  vec3 cd = color * diffuseLightColor * dot(normalize(fragNormal), lightDirection);
  vec3 cs = specularObjectColor * pow(dot(cameraDirection, reflect(lightDirection, fragNormal)), phongExponent);
  outputColor = vec4(ca + cd + cs, 1);
}

// #version 330 core

// in vec2 UV;
// in vec3 Normal;

// uniform sampler2D tex;
// uniform vec3 lightPosition;
// uniform vec3 diffuseLightColor;
// uniform vec3 specularObjectColor;
// uniform vec3 ambientLight;
// uniform float phongExponent;
// uniform vec3 cameraPosition;

// out vec4 colorOut;

// void main()
// {
//   vec3 color = texture(tex, UV).rgb;
//   vec3 fragPosition = vec3(gl_FragCoord.x, gl_FragCoord.y, gl_FragCoord.z);
//   vec3 cameraDirection = normalize(cameraPosition - fragPosition);
//   vec3 lightDirection = normalize(lightPosition - fragPosition);

//   vec3 ca = color * ambientLight;
//   vec3 cd = color * diffuseLightColor * dot(normalize(Normal), lightDirection);
//   vec3 cs = specularObjectColor * pow(dot(cameraDirection, reflect(lightDirection, Normal)), phongExponent);
//   colorOut = vec4(ca + cd + cs, 1);
// }

#version 410

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec2 vertTexCoord;
in vec3 vertNormal;

out vec2 fragTexCoord;
out vec3 fragNormal;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
    fragNormal = vertNormal;
}

// #version 330 core

// in vec3 position;
// in vec3 normal;
// in vec2 uv;

// uniform mat4 mvp;
// uniform vec3 objectPosition;

// out vec3 Color;
// out vec2 UV;
// out vec3 Normal;
// out vec3 fragPosition;

// void main()
// {
//   fragPosition = objectPosition * position;
//   gl_Position = mvp * vec4(position, 1.0);
//   UV = uv;
//   Normal = normal;
// }
