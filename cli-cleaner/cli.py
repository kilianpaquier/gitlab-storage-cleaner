import click

@click.group()
@click.option('--debug', is_flag=True, show_default="false")
@click.pass_context
def cli(ctx, debug: bool):
    """
    Cleaner stands here to help cleans old (and globall all) storage associated to one or multiple gitlab projects.
    """
    # ensure that ctx.obj exists and is a dict (in case `cli()` is called
    # by means other than the `if` block below)
    ctx.ensure_object(dict)

    ctx.obj["DEBUG"] = debug