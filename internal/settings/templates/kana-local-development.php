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
	if ( ! getenv('IS_KANA_ENVIRONMENT') === true ) {
		return;
	}

	if ( ! is_admin() || is_user_logged_in() ) {
		return;
	}

	if ( ! function_exists( 'wp_get_environment_type' ) || 'local' !== wp_get_environment_type() ) {
		return;
	}

	$user = wp_signon(
		array(
			'user_login'    => 'admin',
			'user_password' => 'password',
			'remember'      => true,
		)
	);

	if ( ! is_wp_error( $user ) && isset( $_SERVER['REQUEST_URI'] ) ) {
		wp_safe_redirect( $_SERVER['REQUEST_URI'] );
		exit();
	}
}

add_action( 'init', '\KanaCLI\login_to_admin' );
