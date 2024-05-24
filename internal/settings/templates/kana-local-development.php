<?php
/**
 * Plugin Name: Kana Development Addons
 * Plugin URI: https://github.com/ChrisWiegman/kana
 * Description: Various tweaks to ensure local development is as seamless as possible.
 * Author: Chris Wiegman
 * Version: {{ .Version }}
 *
 * @package KanaCLI
 * @version {{ .Version }}
 **/

 namespace KanaCLI;

/*
 * Disable WordPress updates, checks and emails.
 */
add_filter( 'auto_update_core', '__return_false', 9999 );
add_filter( 'auto_update_plugin', '__return_false', 9999 );
add_filter( 'auto_update_theme', '__return_false', 9999 );
add_filter( 'auto_update_translation', '__return_false', 9999 );
add_filter( 'auto_core_update_send_email', '__return_false', 9999 );
add_filter( 'send_core_update_notification_email', '__return_false', 9999 );
remove_action( 'admin_init', '_maybe_update_core' );
remove_action( 'admin_init', '_maybe_update_plugins' );
remove_action( 'admin_init', '_maybe_update_themes' );

// Set Jetpack to offline mode for easier development.
add_filter( 'jetpack_offline_mode', '__return_true' );

/**
 * Use Mailpit to capture emails from the WordPress site.
 *
 * @param PHPMailer $phpmailer The PHPMailer instance (passed by reference).
 */
function action_phpmailer_init( $phpmailer ) {
	$phpmailer->isSMTP();
	$phpmailer->Host = 'kana-{{ .SiteName }}-mailpit';
	$phpmailer->Port = 1025;
}

add_action( 'phpmailer_init', '\KanaCLI\action_phpmailer_init' );

/**
 * Login to the WordPress admin automatically when visiting a WordPress admin URL.
 */
function login_to_admin() {
	if ( ! getenv('IS_KANA_ENVIRONMENT') === true
		|| ! is_admin()
		|| is_user_logged_in() ) {
		return;
	}

	$kana_error = '<p>Kana could not find a valid admin user to login to your site.</p>';


	$args = array(
		'role'    => 'administrator',
		'orderby' => 'id',
		'order'   => 'ASC',
		'number'  => '2',
	);

	$users = get_users( $args );

	if ( empty( $users ) ) {
		wp_die( $kana_error, 200 );
	}

	$user = $users[0];

	if (! is_wp_error( $user)) {
		wp_set_current_user( $user->ID, $user->user_login );
		wp_set_auth_cookie( $user->ID );

		if ( isset( $_SERVER['REQUEST_URI'] ) ) {
			wp_safe_redirect( $_SERVER['REQUEST_URI'] );
			exit();
		}
	}

	if ( is_wp_error( $user ) ) {
		wp_die( $user->get_error_message() . $kana_error, 200 );
	}
}

add_action( 'set_current_user', '\KanaCLI\login_to_admin' );
