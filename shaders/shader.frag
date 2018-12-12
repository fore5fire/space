#version 410

uniform sampler2D tex;
uniform mat4 camera;

in vec2 fragTexCoord;
in vec3 fragNormal;

out vec4 outputColor;

void main() {

  vec3 fragCamera = vec3(camera[3][0], camera[3][1], camera[3][2]);
  
  vec3 lightPosition = vec3(0,-10,0);
  vec3 diffuseLightColor = vec3(0.8,0.8,0.8);
  vec3 specularObjectColor = vec3(0.8,0.8,0.8);
  vec3 ambientLight = vec3(0.2,0.2,0.2);
  float phongExponent = 1;

  vec3 color = texture(tex, vec2(fragTexCoord.x, 1.0-fragTexCoord.y)).rgb;
  vec3 fragPosition = vec3(gl_FragCoord.x, gl_FragCoord.y, gl_FragCoord.z);
  vec3 cameraDirection = normalize(fragCamera - fragPosition);
  vec3 lightDirection = normalize(lightPosition - fragPosition);

  vec3 ca = color * ambientLight;
  vec3 cd = color * diffuseLightColor * dot(normalize(fragNormal), lightDirection);
  vec3 cs = specularObjectColor * pow(dot(cameraDirection, reflect(lightDirection, fragNormal)), phongExponent);
  outputColor = vec4(ca + cd + cs, 1);
}
