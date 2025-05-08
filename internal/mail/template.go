package mail

func EmailTemplate(name string) string {
	return `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html dir="ltr" lang="es">
  <head>
    <meta content="text/html; charset=UTF-8" http-equiv="Content-Type" />
    <meta name="x-apple-disable-message-reformatting" />
  </head>
  <body
    style='background-color:rgb(240,244,248);font-family:ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";padding-top:40px;padding-bottom:40px'>
    <!--$-->
    <div
      style="display:none;overflow:hidden;line-height:1px;opacity:0;max-height:0;max-width:0">
      Bienvenido a Altheia EHR - ¡Completa tu configuración!
    </div>
    <table
      align="center"
      width="100%"
      border="0"
      cellpadding="0"
      cellspacing="0"
      role="presentation"
      style="background-color:rgb(255,255,255);border-radius:12px;margin-left:auto;margin-right:auto;padding:0px;max-width:600px;box-shadow:var(--tw-ring-offset-shadow, 0 0 #0000), var(--tw-ring-shadow, 0 0 #0000), 0 1px 2px 0 rgb(0,0,0,0.05);overflow:hidden">
      <tbody>
        <tr style="width:100%">
          <td>
            <table
              align="center"
              width="100%"
              border="0"
              cellpadding="0"
              cellspacing="0"
              role="presentation"
              style="background-color:rgb(30,64,175);padding:32px;text-align:center">
              <tbody>
                <tr>
                  <td>
                    <h1
                      style="font-size:32px;font-weight:700;color:rgb(255,255,255);margin:0px">
                      Altheia EHR
                    </h1>
                    <p
                      style="font-size:18px;color:rgb(255,255,255);opacity:0.9;margin-top:8px;margin-bottom:0px;line-height:24px">
                      Tu plataforma de salud digital
                    </p>
                  </td>
                </tr>
              </tbody>
            </table>
            <table
              align="center"
              width="100%"
              border="0"
              cellpadding="0"
              cellspacing="0"
              role="presentation"
              style="padding-left:32px;padding-right:32px;padding-top:40px;padding-bottom:32px">
              <tbody>
                <tr>
                  <td>
                    <h1
                      style="font-size:24px;font-weight:700;color:rgb(30,58,138);margin-bottom:24px;text-align:center">
                      ¡Bienvenido,
                      <span style="color:rgb(59,130,246)">` + name + `</span
                      >!
                    </h1>
                    <p
                      style="font-size:16px;line-height:26px;color:rgb(51,65,85);margin-bottom:24px;margin-top:16px">
                      Nos complace darte la bienvenida a
                      <span style="font-weight:700">Altheia EHR</span>, tu nueva
                      plataforma de Registro Electrónico de Salud. Has tomado
                      una excelente decisión para transformar la gestión de
                      información médica en tu práctica profesional.
                    </p>
                    <table
                      align="center"
                      width="100%"
                      border="0"
                      cellpadding="0"
                      cellspacing="0"
                      role="presentation"
                      style="background-color:rgb(248,250,252);padding:24px;border-radius:8px;margin-bottom:32px;border-left-width:4px;border-color:rgb(59,130,246)">
                      <tbody>
                        <tr>
                          <td>
                            <p
                              style="font-size:18px;font-weight:700;color:rgb(30,58,138);margin-bottom:16px;line-height:24px;margin-top:16px">
                              Con Altheia EHR podrás:
                            </p>
                            <table
                              align="center"
                              width="100%"
                              border="0"
                              cellpadding="0"
                              cellspacing="0"
                              role="presentation"
                              style="margin-bottom:12px">
                              <tbody style="width:100%">
                                <tr style="width:100%">
                                  <td
                                    data-id="__react-email-column"
                                    style="width:24px;padding-right:8px;vertical-align:top">
                                    <p
                                      style="font-size:16px;color:rgb(59,130,246);font-weight:700;margin:0px;line-height:24px;margin-bottom:0px;margin-top:0px;margin-left:0px;margin-right:0px">
                                      ✓
                                    </p>
                                  </td>
                                  <td data-id="__react-email-column">
                                    <p
                                      style="font-size:16px;line-height:24px;color:rgb(51,65,85);margin:0px;margin-bottom:0px;margin-top:0px;margin-left:0px;margin-right:0px">
                                      <span style="font-weight:500"
                                        >Gestionar historias clínicas</span
                                      >
                                      de forma segura y eficiente
                                    </p>
                                  </td>
                                </tr>
                              </tbody>
                            </table>
                            <table
                              align="center"
                              width="100%"
                              border="0"
                              cellpadding="0"
                              cellspacing="0"
                              role="presentation"
                              style="margin-bottom:12px">
                              <tbody style="width:100%">
                                <tr style="width:100%">
                                  <td
                                    data-id="__react-email-column"
                                    style="width:24px;padding-right:8px;vertical-align:top">
                                    <p
                                      style="font-size:16px;color:rgb(59,130,246);font-weight:700;margin:0px;line-height:24px;margin-bottom:0px;margin-top:0px;margin-left:0px;margin-right:0px">
                                      ✓
                                    </p>
                                  </td>
                                  <td data-id="__react-email-column">
                                    <p
                                      style="font-size:16px;line-height:24px;color:rgb(51,65,85);margin:0px;margin-bottom:0px;margin-top:0px;margin-left:0px;margin-right:0px">
                                      <span style="font-weight:500"
                                        >Agendar citas</span
                                      >
                                      y enviar recordatorios automáticos
                                    </p>
                                  </td>
                                </tr>
                              </tbody>
                            </table>
                            <table
                              align="center"
                              width="100%"
                              border="0"
                              cellpadding="0"
                              cellspacing="0"
                              role="presentation"
                              style="margin-bottom:12px">
                              <tbody style="width:100%">
                                <tr style="width:100%">
                                  <td
                                    data-id="__react-email-column"
                                    style="width:24px;padding-right:8px;vertical-align:top">
                                    <p
                                      style="font-size:16px;color:rgb(59,130,246);font-weight:700;margin:0px;line-height:24px;margin-bottom:0px;margin-top:0px;margin-left:0px;margin-right:0px">
                                      ✓
                                    </p>
                                  </td>
                                  <td data-id="__react-email-column">
                                    <p
                                      style="font-size:16px;line-height:24px;color:rgb(51,65,85);margin:0px;margin-bottom:0px;margin-top:0px;margin-left:0px;margin-right:0px">
                                      <span style="font-weight:500"
                                        >Acceder a información médica</span
                                      >
                                      desde cualquier dispositivo
                                    </p>
                                  </td>
                                </tr>
                              </tbody>
                            </table>
                            <table
                              align="center"
                              width="100%"
                              border="0"
                              cellpadding="0"
                              cellspacing="0"
                              role="presentation">
                              <tbody style="width:100%">
                                <tr style="width:100%">
                                  <td
                                    data-id="__react-email-column"
                                    style="width:24px;padding-right:8px;vertical-align:top">
                                    <p
                                      style="font-size:16px;color:rgb(59,130,246);font-weight:700;margin:0px;line-height:24px;margin-bottom:0px;margin-top:0px;margin-left:0px;margin-right:0px">
                                      ✓
                                    </p>
                                  </td>
                                  <td data-id="__react-email-column">
                                    <p
                                      style="font-size:16px;line-height:24px;color:rgb(51,65,85);margin:0px;margin-bottom:0px;margin-top:0px;margin-left:0px;margin-right:0px">
                                      <span style="font-weight:500"
                                        >Mejorar la comunicación</span
                                      >
                                      con tus pacientes
                                    </p>
                                  </td>
                                </tr>
                              </tbody>
                            </table>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <table
                      align="center"
                      width="100%"
                      border="0"
                      cellpadding="0"
                      cellspacing="0"
                      role="presentation"
                      style="background-color:rgb(239,246,255);padding:32px;border-radius:12px;margin-bottom:32px;text-align:center;border-width:1px;border-color:rgb(219,234,254)">
                      <tbody>
                        <tr>
                          <td>
                            <p
                              style="font-size:18px;line-height:28px;color:rgb(30,58,138);font-weight:500;margin-bottom:20px;margin-top:16px">
                              ¡Tu cuenta está casi lista!
                            </p>
                            <p
                              style="font-size:16px;line-height:24px;color:rgb(51,65,85);margin-bottom:24px;margin-top:16px">
                              Para comenzar a disfrutar de todos los beneficios,
                              completa la configuración de tu perfil:
                            </p>
                            <a
                              class="hover:bg-[#1d4ed8]"
                              href="https://altheia-ehr.com/configuracion"
                              style="background-color:rgb(37,99,235);color:rgb(255,255,255);font-weight:700;padding-top:14px;padding-bottom:14px;padding-left:36px;padding-right:36px;border-radius:6px;font-size:16px;text-decoration-line:none;text-align:center;display:inline-block;box-sizing:border-box;box-shadow:var(--tw-ring-offset-shadow, 0 0 #0000), var(--tw-ring-shadow, 0 0 #0000), 0 1px 2px 0 rgb(0,0,0,0.05);line-height:100%;text-decoration:none;max-width:100%;mso-padding-alt:0px;padding:14px 36px 14px 36px"
                              target="_blank"
                              ><span
                                ><!--[if mso]><i style="mso-font-width:450%;mso-text-raise:21" hidden>&#8202;&#8202;&#8202;&#8202;</i><![endif]--></span
                              ><span
                                style="max-width:100%;display:inline-block;line-height:120%;mso-padding-alt:0px;mso-text-raise:10.5px"
                                >Completar Mi Perfil</span
                              ><span
                                ><!--[if mso]><i style="mso-font-width:450%" hidden>&#8202;&#8202;&#8202;&#8202;&#8203;</i><![endif]--></span
                              ></a
                            >
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <table
                      align="center"
                      width="100%"
                      border="0"
                      cellpadding="0"
                      cellspacing="0"
                      role="presentation"
                      style="border-top-width:1px;border-color:rgb(226,232,240);padding-top:24px;margin-top:32px">
                      <tbody>
                        <tr>
                          <td>
                            <p
                              style="font-size:16px;line-height:24px;color:rgb(71,85,105);margin-bottom:16px;margin-top:16px">
                              ¿Necesitas ayuda para comenzar? Nuestro equipo de
                              soporte está listo para asistirte en cada paso del
                              camino.
                            </p>
                            <table
                              align="center"
                              width="100%"
                              border="0"
                              cellpadding="0"
                              cellspacing="0"
                              role="presentation">
                              <tbody style="width:100%">
                                <tr style="width:100%">
                                  <td
                                    data-id="__react-email-column"
                                    style="width:50%;padding-right:8px">
                                    <a
                                      href="https://altheia-ehr.com/ayuda"
                                      style="background-color:rgb(248,250,252);color:rgb(51,65,85);font-weight:500;padding-top:12px;padding-bottom:12px;padding-left:16px;padding-right:16px;border-radius:6px;font-size:14px;text-decoration-line:none;text-align:center;display:block;width:100%;border-width:1px;border-color:rgb(226,232,240);box-sizing:border-box;line-height:100%;text-decoration:none;max-width:100%;mso-padding-alt:0px;padding:12px 16px 12px 16px"
                                      target="_blank"
                                      ><span
                                        ><!--[if mso]><i style="mso-font-width:400%;mso-text-raise:18" hidden>&#8202;&#8202;</i><![endif]--></span
                                      ><span
                                        style="max-width:100%;display:inline-block;line-height:120%;mso-padding-alt:0px;mso-text-raise:9px"
                                        >Centro de Ayuda</span
                                      ><span
                                        ><!--[if mso]><i style="mso-font-width:400%" hidden>&#8202;&#8202;&#8203;</i><![endif]--></span
                                      ></a
                                    >
                                  </td>
                                  <td
                                    data-id="__react-email-column"
                                    style="width:50%;padding-left:8px">
                                    <a
                                      href="https://altheia-ehr.com/contacto"
                                      style="background-color:rgb(248,250,252);color:rgb(51,65,85);font-weight:500;padding-top:12px;padding-bottom:12px;padding-left:16px;padding-right:16px;border-radius:6px;font-size:14px;text-decoration-line:none;text-align:center;display:block;width:100%;border-width:1px;border-color:rgb(226,232,240);box-sizing:border-box;line-height:100%;text-decoration:none;max-width:100%;mso-padding-alt:0px;padding:12px 16px 12px 16px"
                                      target="_blank"
                                      ><span
                                        ><!--[if mso]><i style="mso-font-width:400%;mso-text-raise:18" hidden>&#8202;&#8202;</i><![endif]--></span
                                      ><span
                                        style="max-width:100%;display:inline-block;line-height:120%;mso-padding-alt:0px;mso-text-raise:9px"
                                        >Contactar Soporte</span
                                      ><span
                                        ><!--[if mso]><i style="mso-font-width:400%" hidden>&#8202;&#8202;&#8203;</i><![endif]--></span
                                      ></a
                                    >
                                  </td>
                                </tr>
                              </tbody>
                            </table>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <p
                      style="font-size:16px;line-height:24px;color:rgb(71,85,105);margin-top:32px;margin-bottom:8px">
                      ¡Te deseamos mucho éxito en tu experiencia con Altheia
                      EHR!
                    </p>
                    <p
                      style="font-size:16px;line-height:24px;color:rgb(71,85,105);font-weight:700;margin-bottom:0px;margin-top:16px">
                      El Equipo de Altheia EHR
                    </p>
                  </td>
                </tr>
              </tbody>
            </table>
            <table
              align="center"
              width="100%"
              border="0"
              cellpadding="0"
              cellspacing="0"
              role="presentation"
              style="background-image:linear-gradient(to right, var(--tw-gradient-stops));padding:32px;text-align:center;color:rgb(255,255,255)">
              <tbody>
                <tr>
                  <td>
                    <p
                      style="font-size:14px;opacity:0.9;margin:0px;margin-bottom:0px;line-height:24px;margin-top:0px;margin-left:0px;margin-right:0px">
                      © 2025 Altheia EHR. Todos los derechos reservados.
                    </p>
                    <p
                      style="font-size:14px;opacity:0.9;margin:0px;margin-bottom:0px;line-height:24px;margin-top:0px;margin-left:0px;margin-right:0px">
                      Av. Principal #123, Bogotá, Colombia
                    </p>
                    <p
                      style="font-size:14px;margin:0px;line-height:24px;margin-bottom:0px;margin-top:0px;margin-left:0px;margin-right:0px">
                      <a
                        href="https://altheia-ehr.com/preferencias"
                        style="color:rgb(255,255,255);text-decoration-line:underline;opacity:0.9"
                        >Preferencias de correo</a
                      >
                      •
                      <a
                        href="https://altheia-ehr.com/terminos"
                        style="color:rgb(255,255,255);text-decoration-line:underline;opacity:0.9"
                        >Términos de uso</a
                      >
                      •
                      <a
                        href="https://altheia-ehr.com/privacidad"
                        style="color:rgb(255,255,255);text-decoration-line:underline;opacity:0.9"
                        >Privacidad</a
                      >
                    </p>
                  </td>
                </tr>
              </tbody>
            </table>
          </td>
        </tr>
      </tbody>
    </table>
    <!--7--><!--/$-->
  </body>
</html>
`
}
