import cli
import click
import gitlab

@cli.cli.command()
@click.option(
    "--token", envvar="GITLAB_TOKEN", type=str, show_default="GITLAB_TOKEN",
    help="gitlab read/write token to execute operations") 
@click.option(
    "--server", envvar="CI_SERVER_HOST", type=str, show_default="CI_SERVER_HOST",
    help="gitlab endpoint url")
@click.option(
    "--paths", envvar="CI_PROJECT_NAME", type=str, show_default="CI_PROJECT_NAME", multiple=True,
    help="project names to clean")
@click.option(
    "--dry-run", type=bool, show_default="true",
    help="project names to clean")
@click.option(
    "--threshold", type=str, show_default="CI_PROJECT_NAME",
    help="project names to clean")
@click.pass_context
def clean(_, token: str, url: str, projects: list[str]):
    """
    Retrieves storage associated to input project names and cleans old and unnecessary storage
    """

    gl = gitlab.Gitlab(url, private_token=token, pagination="keyset", order_by="id", per_page=100)

    projects = gl.projects.list()
    for project in projects:
        click.echo(project)

    click.echo("hey !")