import cli
import click

__version__ = "v0.0.0"

@cli.cli.command()
@click.pass_context
def version(_):
    """
    Shows current cleaner version
    """
    click.echo(__version__)